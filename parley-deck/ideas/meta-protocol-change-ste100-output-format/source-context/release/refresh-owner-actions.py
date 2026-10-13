from pathlib import Path
import datetime, hashlib, json, shlex, subprocess, urllib.request

root = Path(__file__).resolve().parent
inventory = json.loads((root / 'owner-source-inventory.json').read_text())
with urllib.request.urlopen('https://registry.npmjs.org/parley-deck-skill', timeout=30) as response:
    registry = json.load(response)
core = json.loads(subprocess.check_output(['parley', 'protocol', 'status', '--json'], text=True))
rows = []
body = ['# Pending owner actions', '',
        'Run these in your attended terminal. Only versions still absent at the recorded check time are listed. No agent has published npm or the global core.', '']
for item in inventory:
    version = item['version']
    archive = item['archives'][0]
    for source in [*item['archives'], item['core']]:
        assert hashlib.sha256(Path(source['path']).read_bytes()).hexdigest() == source['sha256'], source['path']
    npm_missing = version not in registry['versions']
    core_missing = version not in core['installed']
    commands = []
    if npm_missing:
        commands.append('npm publish ' + shlex.quote(archive['path']) + ' --access public')
    if core_missing:
        commands.append('parley protocol publish --version ' + version + ' --from ' + shlex.quote(item['core']['path']))
    row = {'version': version, 'npm_unpublished': npm_missing, 'core_unpublished': core_missing,
           'archive': archive, 'core': item['core'], 'commands': commands}
    rows.append(row)
    if commands:
        body += ['## ' + version, '']
        if npm_missing:
            body += ['npm tarball SHA256: `' + archive['sha256'] + '`', '']
        if core_missing:
            body += ['Staged core SHA256: `' + item['core']['sha256'] + '`', '']
        body += ['```sh', *commands, '```', '']
at = datetime.datetime.now(datetime.timezone.utc).isoformat()
result = {'at': at, 'npm_registry': 'https://registry.npmjs.org/parley-deck-skill',
          'npm_latest': registry.get('dist-tags', {}).get('latest'), 'installed_cores': core['installed'],
          'sources_hash_verified': True, 'versions': rows}
(root / 'owner-publication-status.json').write_text(json.dumps(result, indent=2) + '\n')
body += ['Checked at ' + at + '. npm latest: `' + str(result['npm_latest']) + '`; installed cores: ' + ', '.join(core['installed']) + '.', '']
(root / 'owner-actions-current.md').write_text('\n'.join(body))
print(json.dumps({'at': at, 'npm_latest': result['npm_latest'], 'installed_cores': core['installed'],
                  'pending_commands': sum(len(x['commands']) for x in rows)}, indent=2))
