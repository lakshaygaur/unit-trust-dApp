const router = require('express').Router()
const {invokeChaincode} = require('../app/invoke-transaction')
const {queryChaincode} = require('../app/query')
const config = require('../config.json')
const helper = require('../app/helper')

router.post('/applyAccount',async function(req,res){
    console.log(req.body, req.orgname)
    try {
        let peers = await helper.getPeers(req.orgname)
        let args = [] // name, type
        // args.push()
        args.push(req.username)
        args.push(req.body.type)

        console.log(args)
        let invokeRes = await invokeChaincode(peers, config.channelName, config.chaincodeName, 'applyAccount', args, req.username, req.orgname)
        console.log(invokeRes)

        return res.send({
            status: 'ok',
            msg: 'applied for account successfully'
        })
    } catch (error) {
        return res.send({
            status: 'failed',
            msg: error.message
        })
    }
})


router.post('/sellFund',async function(req,res){
    console.log(req.body, req.orgname)
    try {
        let peers = await helper.getPeers(req.orgname)
        let args = [] // type, value, validFrom, validTo, owner
        // args.push()
        args.push(req.body.fundId)
        args.push(req.body.sellingTo)
        console.log(args)
        let invokeRes = await invokeChaincode(peers, config.channelName, config.chaincodeName, 'sellFund', args, req.username, req.orgname)
        console.log(invokeRes)

        return res.send({
            status: 'ok',
            msg: 'fund sold successfully',
            txnIDd : invokeRes
        })
    } catch (error) {
        return res.send({
            status: 'failed',
            msg: error.message
        })
    }
})


module.exports = router