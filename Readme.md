# GO-Online user service proxy

## Introduction

This service proxy request from user to user container.

## Info

### PORT

    This service listening on port `8081`.

### 依赖的环境变量

- CONSUL_ADDRESS consul地址，默认为localhost
- CONSUL_PORT consul端口，默认为8500
- DOMAIN_NAME 域名，在解析时使用，默认为`localhost`
