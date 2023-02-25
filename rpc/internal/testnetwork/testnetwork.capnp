@0xbcea0965c2a55c5b;
# This schema defines the concrete format of third-party handoff
# related data types used by the test network.

using Go = import "/go.capnp";
$Go.package("testnetwork");
$Go.import("capnproto.org/go/capnp/v3/rpc/internal/testnetwork");

struct PeerAndNonce {
  # A pair of peer ID and a nonce. This is the format for all
  # three of ProvisionId, RecipientId, and ThirdPartyCapId,
  # though which peer the id refers to differs.

  peerId @0 :UInt64;
  nonce @1 :UInt64;
}

using ProvisionId = PeerAndNonce;
# peerId is that of the peer that sends the provide.

using RecipientId = PeerAndNonce;
# peerId is that of the peer that sends the accept.

using ThirdPartyCapId = PeerAndNonce;
# peerId is that of the peer that hosts the capability.
