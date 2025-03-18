# Cap'n Proto schema for flinkcoin.core
@0x9a7b8c5d4e3f2a1b; # Unique file ID

using Go = import "/go.capnp";
$Go.package("github.com/flinkcoin/mono/libs/schema/pkg/core");
$Go.import("github.com/flinkcoin/mono/libs/schema/pkg/core");

struct Block {
  # Nested enums
  signatureMode @0 :SignatureMode;
  enum SignatureMode {
    oneOfOne @0;
    twoOfThree @1;
    threeOfThree @2;
    threeOfFive @3;
    fiveOfFive @4;
  }

  blockType @1 :BlockType;
  enum BlockType {
    create @0;
    send @1;
    receive @2;
    update @3;
  }

  # Nested struct for PublicKeys
  publicKeys @2 :PublicKeys;
  struct PublicKeys {
    signatureMode @0 :SignatureMode;
    publicKey @1 :List(Data); # repeated bytes becomes List(Data)
  }

  # Body struct
  body @3 :Body;
  struct Body {
    version @0 :Int32;
    timestamp @1 :Int64;
    blockType @2 :BlockType;
    previousBlockHash @3 :Data;
    accountId @4 :Data;
    delegatedNodeId @5 :Data;
    balance @6 :Int64;
    amount @7 :Int64;
    sendAccountId @8 :Data;
    receiveBlockHash @9 :Data;
    referenceCode @10 :Data;
    publicKeys @11 :PublicKeys;
  }

  # Signatures struct
  signatures @4 :Signatures;
  struct Signatures {
    signature @0 :List(Data); # repeated bytes becomes List(Data)
  }

  # Hash struct
  blockHash @5 :Hash;
  struct Hash {
    hash @0 :Data;
  }

  # Work struct
  work @6 :Work;
  struct Work {
    work @0 :Data;
  }
}

struct PaymentRequest {
  fromAccountId @0 :Data;
  toAccountId @1 :Data;
  amount @2 :Int64;
  referenceCode @3 :Data;
}

struct FullBlock {
  block @0 :Block;
  next @1 :Data;
}

struct Node {
  nodeId @0 :Data;
  publicKey @1 :Data;
}

struct NodeAddress {
  ip @0 :Text;
  port @1 :Int32;
}