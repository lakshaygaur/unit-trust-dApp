const helper = require('./helper.js')
const config = require('./config.json')
var Account = require('../../models/account.js');
var passport = require('passport');
const mysql = require('mysql')
const db = require('../../webapp/db.js')
require('./config.js')
async function regsiterAdmin(){
    try {
    let secret =  await helper.getRegisteredUser(config.admin,config.org ,true,config.admin_pwd);
    console.log('Admin pwd ', secret)
    let admin = {
        address : config.admin,
        name : config.admin,
        email : config.admin,
        company_name : config.org,
        type: '7'
    }
        await db.query(mysql.format('insert into Account set ? ',[admin]))     
    Account.register(new Account({ username : config.admin, org: config.org }), config.admin_pwd, function(err, account) {
        console.log('here')
        if (err) {
        console.log('error',err)
        }
        return process.exit()
    })
    } catch (error) {
        console.log(error)
    }
}

regsiterAdmin()
