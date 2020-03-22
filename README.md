## Unit Trust Blockchain dApp

A blockchain network for a Unit Trust company called MTCT including a headquarter (HQ), Two branches, Two financial advisers(agent), and Six investors. The HQ and its branches maintain the network. All agree on the type of assets e.g Bond, Equity etc. to create on the network. An agent service two investors. Investors view only assets on the network. The HQ and branches can view all agents, assets, and investors on the network. The HQ and branches all agree before adopting any new policy. Investors can buy or sell their fund back to MTCT.



### User Story
* Agent can apply to join the network.
* Agent can sell fund to investor.
* Investor can view the type of available fund on the network.
* Branch office can onboard agents on the network.
* MTCT can accept or reject application to join the network.
* MTCT can create or destroy any fund from the network.

### Prerequisites:

* [Docker](https://www.docker.com/products/overview) - v1.12 or higher
* [Docker Compose](https://docs.docker.com/compose/overview/) - v1.8 or higher
* [Git client](https://git-scm.com/downloads) - needed for clone commands
* **Node.js** v8.4.0 or higher
* [Download Docker images](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#binaries)

### Setup: 
##### Terminal Window 1
```
cd artifacts
docker-compose up
```
##### Terminal Window 2
```
npm install
node dev/dev_script/channel_script.js
```
Wait for script to bootstrap the network. Then start the application with following command -

` node app.js `
The application by defualt runs on PORT=4000. You can change it in *config.json*

Once you have completed the above setup, you will have provisioned a local network with the following docker container configuration:
* 3 Orgs -
* * Hq (MTCT Headquarter)
* * Branch 1
* * Branch 2
* 2 CAs
* A SOLO orderer
* 3 peers (1 peers per Org)

#### Artifacts
* Crypto material has been generated using the **cryptogen** tool from Hyperledger Fabric and mounted to all peers, the orderering node and CA containers.
* An Orderer genesis block (genesis.block) and channel configuration transaction (mychannel.tx) has been pre generated using the **configtxgen** tool from Hyperledger Fabric and placed within the artifacts folder.


## REST APIs

### Generate Certificate

* Register and enroll new users in Organization - **Hq**, **Branch1** or **Branch2**:
```
curl --location --request POST 'http://localhost:4000/users' \
--header 'Content-Type: application/json' \
--data-raw '{
	"username":"JimFromHq",
	"orgName":"Hq"
}'
```

**OUTPUT:**

```
{
  "success": true,
  "secret": "RaxhMgevgJcm",
  "message": "Jim enrolled Successfully",
  "token": "<put JSON Web Token here>"
}
```

The response contains the success/failure status, an **enrollment Secret** and a **JSON Web Token (JWT)** that is a required string in the Request Headers for subsequent requests.

### MTCT - Create Account

```
curl --location --request POST 'http://localhost:4000/mtct/createAccount' \
--header 'Content-Type: application/json' \
--data-raw '{
	"type": "MTCT"
}'
```
**OUTPUT:**

```
{
    "status": "ok",
    "msg": "account created successfully"
}
```

Please note that the Header **authorization** must contain the JWT returned from the `POST /users` call

### MTCT - Create Fund

```
curl --location --request POST 'http://localhost:4000/mtct/createFund' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <jwt token>' \
--data-raw '{
	"type": "Bond",
	"value": "$1300",
	"validFrom": "03-02-2019",
	"validTo": "03-02-2021"
}'
```
**OUTPUT:**

```
{
    "status": "ok",
    "msg": "fund created successfully",
    "txnIDd": "7f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191"
}
```

### MTCT - Read Fund with Transaction History

```
curl --location --request GET 'http://localhost:4000/mtct/readFund/7f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191' \
--header 'Authorization: Bearer <jwt token>'
```
**OUTPUT:**

```
{
  "fund": {
    "fundId": "\u0000FUND\u00007f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191\u0000",
    "type": "Bond",
    "value": "$1300",
    "validFrom": "03-02-2019",
    "validTo": "03-02-2021",
    "owner": "InvestorBranch1"
  },
  "txnHistory": [
    {
      "txnId": "\u0000TRANSACTION_HISTORY\u00007f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191\u000020a8bf0f53ca31df9b9f2df496da4c470928a7a56b719fe90ddad54e67653634\u0000",
      "status": "Fund: Bond sold to InvestorBranch1",
      "timestamp": "seconds:1584867118 nanos:902000000 "
    },
    {
      "txnId": "\u0000TRANSACTION_HISTORY\u00007f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191\u00007f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191\u0000",
      "status": "Fund: Bond created by JimFromHq",
      "timestamp": "seconds:1584865813 nanos:546000000 "
    }
  ]
}
```
### MTCT - Approve agent's account
```
curl --location --request POST 'http://localhost:4000/mtct/approveAccount' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <Jwt token>' \
--data-raw '{
	"agentId": "JimAgent"
}'
```
**OUTPUT:**

```
{
    "status": "ok",
    "msg": "account approved successfully",
    "txnIDd": "b6dbcafbc641b3cd097e1b9e4c4d98195834cb89b70c64682046f82638abfeb8"
}
```

### MTCT - Approve agent's account
```
curl --location --request POST 'http://localhost:4000/mtct/deleteFund' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <jwt token>' \
--data-raw '{
	"fundId": "7f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191"
}'
```
**OUTPUT:**

```
{
    "status": "ok",
    "msg": "fund deleted successfully",
    "txnIDd": "17ee70e4f781ba274dab9e76ea810cae8f59e6f6c91262413c51ff766bc8c714"
}
```
### Investor - Read all funds available in network
```
curl --location --request GET 'http://localhost:4000/investor/readAllFunds' \
--header 'Authorization: Bearer <JWT TOKEN of investor>'
```
**OUTPUT:**

```
{
  "funds": [
    {
      "fundId": "\u0000FUND\u00004391875f6830151af64e17bc6f6f587ab474bf7fe86302f990526d68af5e855c\u0000",
      "type": "Bond",
      "value": "$1300",
      "validFrom": "03-02-2019",
      "validTo": "03-02-2021",
      "owner": "JimFromHq"
    },
    {
      "fundId": "\u0000FUND\u00007f1dbcc6b242a1327083e5dbcf49ce2ee99ef9460c375eb9e065c1b4602da191\u0000",
      "type": "Equity",
      "value": "$100",
      "validFrom": "03-06-2020",
      "validTo": "03-05-2021",
      "owner": "TimFromHq"
    }
  ]
}
```

### Agent - Apply for account
**NOTE**: Agent first needs to be register and enrolled from either of the Branch orgs. Only then agent can apply to join the network and get his account approved from MTCT(HQ).
```
curl --location --request POST 'http://localhost:4000/agent/applyAccount' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <Jwt token of agent>' \
--data-raw '{
	"type": "AGENT"
}'
```
**OUTPUT:**
```
{
    "status": "ok",
    "msg": "applied for account successfully"
}
```

### Agent - Sell funds 
**NOTE**: This request will result in error if the agent's account has not been approved from MTCT first.
```
curl --location --request POST 'http://localhost:4000/agent/sellFund' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <Jwt token of agent>' \
--data-raw '{
	"fundId": "6ffe82f0c4dba64d4d3a8c8609861f24c4c083486e96ca63772fbb3c134b0f6c",
	"sellingTo": "InvestorBranch1"
}'
```

## To Do
- [x] Agent can apply to join the network.
- [x] Agent can sell fund to investor.
- [x] Investor can view the type of available fund on the network.
- [x] Branch office can onboard agents on the network.
- [x] MTCT can accept or reject application to join the network.
- [x] MTCT can create or destroy any fund from the network.
- [ ] Investor can buy or sell fund.
- [ ] Branch office can create a private route to transact
- [ ] MTCT can reject application of an agent to join the network

### Clean the network

The network will still be running at this point. Before starting the network manually again, here are the commands which cleans the containers and artifacts.

```
cd artifacts/
docker-compose down
cd -
./dev/dev_script/deleteImages.sh
```


