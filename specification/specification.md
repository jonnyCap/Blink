# Fast TCP Pub/Sub Protocol Specification

## Overview

This specification defines a minimalist, high-performance protocol for a publish/subscribe queue system over raw TCP. The protocol is designed for low-latency environments and emphasizes fixed-size control messages, efficient parsing, and authentication using JWTs.

## Protocol Principles

- Connection-oriented (TCP)
- Fixed-size control packets for fast parsing
- Minimal framing overhead
- Stateless authentication using JWT in each request
- Per-queue authorization via key in JWT
- Topic IDs (hashed uint32) instead of string-based topic names for efficiency
- Support for dynamic key rotation to revoke access

---

## Connection Lifecycle

1. **Connect**: Client establishes a persistent TCP connection to the server.
2. **CREATE**: Client requests the creation of a topic/queue, setting a per-queue access key.
3. **SUBSCRIBE**: Client subscribes to an existing topic using a JWT with correct permissions.
4. **PUBLISH**: Client publishes messages to a topic using a JWT with correct permissions.
5. **MESSAGE**: Server delivers messages to subscribed clients.
6. **UNSUBSCRIBE (optional)**
7. **ROTATE_KEY**: Client with appropriate permissions requests key rotation.
8. **KEY_UPDATE (server → client)**: Server distributes a new access key to valid clients after rotation.

---

## Packet Structure

Each packet begins with a 1-byte type identifier. Most commands use fixed sizes for predictability.

### Common Fields

- `Type` (1 byte): Command identifier
- `TopicID` (4 bytes): uint32, hashed topic identifier (used after CREATE)

---

## Packet Types and Formats

### 0x01: CREATE

Create a topic (queue) and register an access key.

| Field        | Size    | Description                 |
| ------------ | ------- | --------------------------- |
| Type         | 1 byte  | 0x01                        |
| JWT Length   | 1 byte  | N (<= 255)                  |
| JWT          | N bytes | JWT containing `queue_key`  |
| Topic Length | 1 byte  | T (<= 255)                  |
| Topic Name   | T bytes | UTF-8 topic string          |
| Flags        | 1 byte  | Optional: durable/exclusive |

### 0x02: SUBSCRIBE

Subscribe to a topic.

| Field      | Size    | Description                  |
| ---------- | ------- | ---------------------------- |
| Type       | 1 byte  | 0x02                         |
| JWT Length | 1 byte  | N (<= 255)                   |
| JWT        | N bytes | JWT with `queue_key` + perms |
| TopicID    | 4 bytes | uint32 hash                  |

### 0x03: PUBLISH

Publish a message to a topic.

| Field          | Size    | Description                  |
| -------------- | ------- | ---------------------------- |
| Type           | 1 byte  | 0x03                         |
| JWT Length     | 1 byte  | N (<= 255)                   |
| JWT            | N bytes | JWT with `queue_key` + perms |
| TopicID        | 4 bytes | uint32                       |
| Payload Length | 4 bytes | uint32                       |
| Payload        | N bytes | Raw payload                  |

### 0x04: MESSAGE

Server-to-client delivery of a published message.

| Field          | Size    | Description |
| -------------- | ------- | ----------- |
| Type           | 1 byte  | 0x04        |
| TopicID        | 4 bytes | uint32      |
| Payload Length | 4 bytes | uint32      |
| Payload        | N bytes | Raw payload |

### 0x05: UNSUBSCRIBE (optional)

| Field      | Size    | Description                  |
| ---------- | ------- | ---------------------------- |
| Type       | 1 byte  | 0x05                         |
| JWT Length | 1 byte  | N (<= 255)                   |
| JWT        | N bytes | JWT with `queue_key` + perms |
| TopicID    | 4 bytes | uint32 hash                  |

### 0x06: KEY_UPDATE (server → client)

Notify authorized clients of a new `queue_key` for a topic.

| Field       | Size    | Description           |
| ----------- | ------- | --------------------- |
| Type        | 1 byte  | 0x06                  |
| TopicID     | 4 bytes | uint32                |
| New Key Len | 1 byte  | N                     |
| New Key     | N bytes | New access key string |

### 0x07: ROTATE_KEY

Client requests the server to rotate the key for a topic. Requires `rotate` permission.

| Field       | Size    | Description                   |
| ----------- | ------- | ----------------------------- |
| Type        | 1 byte  | 0x07                          |
| JWT Length  | 1 byte  | N (<= 255)                    |
| JWT         | N bytes | JWT with `queue_key` + rotate |
| TopicID     | 4 bytes | uint32 hash                   |
| New Key Len | 1 byte  | M                             |
| New Key     | M bytes | New queue access key          |

---

## Authentication

- Server is initialized with a master secret key.
- JWTs are signed using this key.
- During `CREATE`, the client includes a `queue_key` claim in the JWT, which is then stored and bound to the created queue.
- For `PUBLISH`, `SUBSCRIBE`, or `UNSUBSCRIBE`, the JWT must contain:

  - A matching `queue_key` for the topic
  - A `permissions` claim such as `"publish"` or `"subscribe"`

- To rotate a key, clients must use the `ROTATE_KEY` command with a `rotate` permission in the JWT.
- Server verifies:

  1. Signature of the JWT
  2. That the `queue_key` matches the queue's access key
  3. That the user has appropriate permission in the token

### JWT Example (Standard Access)

```json
{
  "sub": "user123",
  "queue_key": "abc123xyz",
  "permissions": ["publish", "subscribe"]
}
```

### JWT Example (Key Rotation Authority)

```json
{
  "sub": "queue_owner",
  "queue_key": "abc123xyz",
  "permissions": ["rotate"]
}
```

## Topic IDs

- Topics are referenced via 32-bit hashed identifiers.
- Servers maintain a mapping from hashed TopicID to string/topic metadata.

---

## Security Notes

- All stateful commands require valid JWTs.
- Per-queue `queue_key` prevents unauthorized access across queues.
- JWTs should be signed using HMAC-SHA256 or similar.
- Key rotation enables client eviction without explicit connection management.
- Consider implementing rate limiting and connection quotas.

---

## Extensibility

- New command types can be added with unused packet codes.
- Flags in CREATE can encode queue types (durable, fanout, TTL).
- Topic name can be replaced with numeric ID entirely if preconfigured.

---

## Example Sequence

1. Client connects
2. Client sends CREATE with `queue_key` in JWT
3. Client sends SUBSCRIBE to topic ID using JWT with same `queue_key`
4. Client sends PUBLISH using JWT with `queue_key` and permission
5. Queue owner sends ROTATE_KEY with `rotate` permission and new key
6. Server sends KEY_UPDATE to all valid subscribers
7. Server sends MESSAGE to all still-authorized clients

---

## Summary

This protocol is designed to be:

- Minimal and fast
- Easy to implement in any language
- Scalable and secure with minimal overhead

Use this as the foundation for a lightweight, embedded message broker or a highly optimized pub/sub layer in distributed systems.
