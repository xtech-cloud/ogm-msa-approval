package config

const defaultYAML string = `
service:
  name: omo.api.msa.approval
  address: :9606
  ttl: 15
  interval: 10
logger:
  level: info
  dir: /var/log/msa/
database:
  lite: true
  timeout: 10
  mysql:
    address: 127.0.0.1:3306
    user: root
    password: mysql@OMO
    db: msa_approval
  sqlite:
    path: /tmp/msa-approval.db
publisher:
- /workflow/make
- /workflow/remove
- /operator/join
- /operator/leave
- /task/submit
- /task/accept
- /task/reject
`
