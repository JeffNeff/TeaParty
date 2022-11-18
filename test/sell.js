var axios = require('axios');

const NKNADDRESS = "10be12ec0133e40e00e13b9c6bc127ded8def6873bf96b67d16ac0721f9dbf81"
const SHIPPINGADDRESS = "0x83403ef2313adFF264a4fE9d5629945D7D6d12C5"


var data = JSON.stringify({
    "currency": "mineonlium",
    "amount": 1000000000000000000,
    "tradeAsset": "ethereum",
    "price": 2000000000000000000,
    "txid": "0x1",
    "locked": false,
    "sellerShippingAddress": SHIPPINGADDRESS,
    "sellerNKNAddress": NKNADDRESS
});

var config = {
    method: 'post',
    url: 'http://0.0.0.0:8080/sell',
    headers: {
        'Content-Type': 'application/json'
    },
    data: data
};

axios(config)
    .then(function (response) {
        console.log(JSON.stringify(response.data));
    })
    .catch(function (error) {
        console.log(error);
    });
