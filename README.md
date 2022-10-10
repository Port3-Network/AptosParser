# Port3-Network parse

## [sync's deploy documentation](./SyncDeploy.md)
## Coin, event flag: 0x1::coin

#### call function:
    Mint: 0x1::managed_coin::mint
    Transfer: 0x1::coin::transfer

#### event type:
    Mint: 0x1::coin::DepositEvent
    Burn: 0x1::coin::WithdrawEvent

#### transfer event: 
    Transfer: Address A -> Address B, amount 100. two events will be visibe
    Address A event: type -> 0x1::coin::WithdrawEvent, amount -> 100
    Address B event: type -> 0x1::coin::DepositEvent, amount -> 100

#### mint event:
    Mint 100 coin to Address A
    Address A event: type -> 0x1::coin::DepositEvent, amount -> 100

#### Conclusion: when a transfer event occurs, Address A will Burn x number of coin and Address B will be Mint with x number of coin

## Token, event flag: 0x3::token

### Aptos references the collection, nft minting on collection
    collection: when you need to create nft, create the collection first, set the collection name. Basically, collection is a combination of <creator + collection name>, it means uniqueness.
    NFT: when you need to create nft, provide the collection name, nft name, basically, nft is combination of <creator + collection name + nft name>, it means uniqueness
#### event type:
    Collection create: 0x3::token::CreateCollectionEvent
    Token create: 0x3::token::CreateTokenDataEvent
    Token Mint: 0x3::token::DepositEvent
    Token burn: 0x3::token::WithdrawEvent
#### call function
    Collection create: 0x3::token::create_collection_script
    Token create: 0x3::token::create_token_script
    Token offer:0x3::token_transfers::offer_script
    Token claim: 0x3::token_transfers::claim_script
    Token direct: 0x3::token::direct_transfer_script
