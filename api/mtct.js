const router = require('express').Router()
const {invokeChaincode} = require('../app/invoke-transaction')
const {queryChaincode} = require('../app/query')
const config = require('../config.json')
const helper = require('../app/helper')

router.post('/createAccount',async function(req,res){
    console.log(req.body, req.orgname)
    try {
        let peers = await helper.getPeers(req.orgname)
        let args = [] // name, type
        // args.push()
        args.push(req.username)
        args.push(req.body.type)

        console.log(args)
        let invokeRes = await invokeChaincode(peers, config.channelName, config.chaincodeName, 'createAccount', args, req.username, req.orgname)
        console.log(invokeRes)

        return res.send({
            status: 'ok',
            msg: 'fund created successfully'
        })
    } catch (error) {
        return res.send({
            status: 'failed',
            msg: error.message
        })
    }
})


router.post('/createFund',async function(req,res){
    console.log(req.body, req.orgname)
    try {
        let peers = await helper.getPeers(req.orgname)
        let args = [] // type, value, validFrom, validTo, owner
        // args.push()
        args.push(req.body.type)
        args.push(req.body.value)
        args.push(req.body.validFrom)
        args.push(req.body.validTo)
        console.log(args)
        let invokeRes = await invokeChaincode(peers, config.channelName, config.chaincodeName, 'createFund', args, req.username, req.orgname)
        console.log(invokeRes)

        return res.send({
            status: 'ok',
            msg: 'fund created successfully',
            txnIDd : invokeRes
        })
    } catch (error) {
        return res.send({
            status: 'failed',
            msg: error.message
        })
    }
})

router.get('/readFund/:fundId',async function(req,res){
    try {
        let peers = await helper.getPeers(req.orgname)
        let args = []
        args.push(req.params.fundId)

        let queryRes = await queryChaincode(peers[0], config.channelName, config.chaincodeName, args, 'readFund', req.username, req.orgname)
        return res.send(JSON.parse(queryRes))
    } catch (error) {
        return res.send({
            status: 'failed',
            msg: error.message
        })
    }
})

module.exports = router