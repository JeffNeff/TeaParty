# TeaParty
This is the public alpha testing of Tea and Party. 

The goal of `alpha1v1` release:
* Begin transfering funds between parties from one chain to another using the proposed escrow system.

Security, fail over, scaling, and data/state management have not been addresed yet. 

**REMINDERS AND ANNOUCEMENTS:**
* This should NOT be connected to any mainnet RPC's.
* This is strictly Proof Of Concept code and still has tons of flaws
* This is set up for debugging, I.E. we are logging and storing lots of things that 
would not, and should not, be stored in a production enviorment. PLEASE DO NOT PUT REAL MONEY INTO THIS.

## Known Issues

* The refund/failover system is not completed. So failed transactions right now == lost transactions 


## Start Tea and Party locally. 

Populate the env. sections with RPC servers. 

Then bring up Tea and Party
```
docker compose up -d
docker compose -f docker-compose.tea.yaml up -d 
```

visit http://localhost:8081 to view the debugging panel. 



## Interacting with Party

Untill `tea` is complete, we interact with `party` via a set of scripts found in the `test/` folder

### create a trade order:

* Update the const at the top of the `sell.js` file and modify the variables at the top to reflect your own NKN and Shipping addresses. 

* Run the `sell.js` file 

```
node sell.js
```

### list the open orders:

```
node list.js
```

### purchase an order:

* Update the const at the top of the `buy.js` file and modify the variables at the top to reflect your own NKN and Shipping addresses. 

```
node buy.js
```
