# homework_test

## Description
A mini client to visit certain ETH endpoint.  
`client` folder contains all the code interacts with ETH.  
`util` folder contains helper function    
`parser.go` contain the core logic how to store address and transaction    
`main.go` set up a gin server for serving request  

## API list
### `get` /current_block
get last recorded block num

### `post` /:address/subscribe
subscribe one user address in system

### `get` /:address/transaction
get all transaction that related to this address