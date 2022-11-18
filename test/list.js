var axios = require('axios');

var config = {
    method: 'get',
    url: 'http://0.0.0.0:8080/listorders',
};

axios(config)
    .then(function (response) {
        console.log(JSON.stringify(response.data));
    })
    .catch(function (error) {
        console.log(error);
    });
