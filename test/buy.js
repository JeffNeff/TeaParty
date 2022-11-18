var axios = require('axios');

const NKNADDRESS = "10be12ec0133e40e00e13b9c6bc127ded8def6873bf96b67d16ac0721f9dbf81"
const SHIPPINGADDRESS = "0x83403ef2313adFF264a4fE9d5629945D7D6d12C5"


var data = JSON.stringify({
  "txid": "0x1",
  "buyerShippingAddress": SHIPPINGADDRESS,
  "buyerNKNAddress": NKNADDRESS
});

var config = {
  method: 'post',
  url: 'http://0.0.0.0:8080/buy',
  headers: { 
    'Content-Type': 'application/json'
  },
  data : data
};

axios(config)
.then(function (response) {
  console.log(JSON.stringify(response.data));
})
.catch(function (error) {
  console.log(error);
});
