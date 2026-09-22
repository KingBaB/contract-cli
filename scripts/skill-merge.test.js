const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const test = require('node:test');

// SKILL_TEST_ROOT also checks the files exported by `skills install`.
const root = process.env.SKILL_TEST_ROOT || path.join(__dirname, '..', 'skills');
const read = name => fs.readFileSync(path.join(root, name, 'SKILL.md'), 'utf8');

test('every shipped skill has unambiguous top-level frontmatter keys', () => {
  for (const entry of fs.readdirSync(root, { withFileTypes: true })) {
    if (!entry.isDirectory() || !fs.existsSync(path.join(root, entry.name, 'SKILL.md'))) continue;
    const frontmatter = read(entry.name).trimStart().match(/^---\r?\n([\s\S]*?)\r?\n---/);
    assert.ok(frontmatter, `${entry.name}: missing frontmatter`);
    const keys = [...frontmatter[1].matchAll(/^([\w-]+):/gm)].map(match => match[1]);
    assert.equal(new Set(keys).size, keys.length, `${entry.name}: duplicate metadata keys: ${keys}`);
    assert.ok(keys.includes('name') && keys.includes('description'), entry.name);
  }
});

test('shared skill gives contract search one destination while preserving review and directory routes', () => {
  const shared = read('contract-cli-shared');
  const decisions = shared.split('## 快速决策')[1].split('\n## ')[0];
  const searchRoutes = decisions.split('\n').filter(line => line.startsWith('- 合同搜索'));
  assert.equal(searchRoutes.length, 1, 'contract search must not route to two skills');
  assert.ok(searchRoutes[0].includes('../contract-cli-contract-search/SKILL.md'));
  for (const name of ['contract-cli-contract-search', 'contract-cli-contract', 'contract-cli-review', 'contract-cli-review-rule', 'contract-cli-employee', 'contract-cli-department', 'contract-cli-rule']) {
    assert.ok(decisions.includes(`../${name}/SKILL.md`), `missing route ${name}`);
    assert.ok(fs.existsSync(path.join(root, name, 'SKILL.md')), `missing destination ${name}`);
  }
});

test('Device quick flow checks and authorizes the same test profile', () => {
  const flow = read('auth').split('### 豆包 / WorkBuddy / Codex Device 授权')[1].split('### `user` 登录')[0];
  const commands = [...flow.matchAll(/^contract-cli auth (status|init|complete) --profile (\S+)/gm)];
  assert.deepEqual(commands.map(match => match[1]), ['status', 'init', 'complete']);
  assert.deepEqual([...new Set(commands.map(match => match[2]))], ['contract-test']);
});

test('shared auth instructions do not override the selected test profile', () => {
  const shared = read('contract-cli-shared');
  const checks = [...shared.matchAll(/auth status --profile (\S+) --as (user|app)/g)];
  assert.ok(checks.length >= 2, 'both identities need authorization checks');
  for (const check of checks) {
    assert.ok(['<profile>', 'contract-test'].includes(check[1]), `unexpected fixed profile ${check[1]}`);
  }
});

test('app authorization examples retain the selected profile', () => {
  const flow = read('auth').split('### `app` 登录')[1].split('## 状态、退出与切换')[0];
  const commands = flow.split('\n').filter(line => line.startsWith('contract-cli auth '));
  assert.ok(commands.length >= 3);
  for (const command of commands) {
    assert.match(command, /--profile <profile>(?:\s|$)/, command);
  }
});

test('test environment selection does not force a different profile name', () => {
  assert.ok(!read('contract-cli-shared').includes('用户指定 test 后，后续所有授权和业务命令均显式传入 `--profile contract-test`'));
});
