# X-Block Restful API

* [Introduction](#introduction)
* [Restful Api List](#restful-api-list)

Restful Api List

| Method | url |
| :---| :---|
| get_gen_blk_time | GET /api/v1/node/generateblocktime |
| get_conn_count | GET /api/v1/node/connectioncount |
| get_blk_txs_by_height | GET /api/v1/block/transactions/height/:height |
| get_blk_by_height | GET /api/v1/block/details/height/:height?raw=0 |
| get_blk_by_hash | GET /api/v1/block/details/hash/:hash?raw=1 |
| get_blk_height | GET /api/v1/block/height |
| get_blk_hash | GET /api/v1/block/hash/:height |
| get_tx | GET /api/v1/transaction/:hash |
| get_balance | GET /api/v1/balance/:addr |
| get_contract_state | GET /api/v1/contract/:hash |
| get_smtcode_evt_txs | GET /api/v1/smartcode/event/transactions/:height |
| get_smtcode_evts | GET /api/v1/smartcode/event/txhash/:hash |
| get_blk_hgt_by_txhash | GET /api/v1/block/height/txhash/:hash |
| get_gasprice | GET /api/v1/gasprice|
| get_allowance | GET /api/v1/allowance/:asset/:from/:to |
| get_unclaimxcg | GET /api/v1/unclaimxcg/:addr |
| post_raw_tx | post /api/v1/transaction?preExec=0 |


## Introduction

This document describes the restful api format for the http/https used in the X-Block.

## Restful Api List

### Response parameters descri

| Field | Type | Description |
| :--- | :--- | :--- |
| Action | string | action name |
| Desc | string | description |
| Error | int64 | error code |
| Result | object | execute result |
| Version | string | version information |


