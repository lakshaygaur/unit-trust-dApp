const router = require('express').Router()
// const {invokeChaincode} = require('../app/invoke-transaction')
const {queryChaincode} = require('../app/query')
const config = require('../config.json')
const helper = require('../app/helper')

router.get('/readAllFunds',async function(req,res){
    try {
        let peers = await helper.getPeers(req.orgname)

        let queryRes = await queryChaincode(peers[0], config.channelName, config.chaincodeName, [], 'readAllFunds', req.username, req.orgname)
        return res.send(JSON.parse(queryRes))
    } catch (error) {
        return res.send({
            status: 'failed',
            msg: error.message
        })
    }
})

module.exports = router