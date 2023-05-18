# shear-plate-sharing

## what is that?
when you use different os in different computer,you may need share shear plate,
so this program is the answer.

## how to use?
### 1、edit `config.yml`
server:
```yaml
Role: server
Port: 7777
Password: xxx
```
client:
```yaml
Host: 192.168.31.174
Port: 7777
Password: xxx
```
tips: do not run client when you pc has been run server

### 2、run
```shell
export CGO_ENABLED=1
go run ./
```

## todo
1、make the code perfect  
2、support copy img