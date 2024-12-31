# Managed DNS Server

A DNS server implementation written in Go from **scratch**,
supporting a management service via a RESTful API.
This project includes custom DNS record handling, database
integration for record persistence, and configuration management.

## Features

- Custom _authorative_ DNS server and client implementation.
- REST API for managing DNS records.
- Supports `A`, `CNAME`, and other common DNS record types.
- Uses PostgreSQL for persistent record storage.
- Implements EDNS (Extension Mechanisms for DNS) with optional support.

# Background

## How DNS Resolution Works

When you type a domain name into your browser (e.g., `www.example.com`), your computer needs to convert that domain name into an
IP address so it can communicate with the correct server. This process is known as **DNS resolution** and involves several steps:

1. **Root Name Servers**  
    The process begins with a request to a **root name server**.
   The root name servers are responsible for directing the query to the appropriate **Top-Level Domain (TLD)** server based
   on the suffix of the domain name (e.g., `.com`, `.org`, `.net`).

2. **TLD Servers**  
    The TLD servers are responsible for handling the domain extensions (such as `.com`, `.org`, etc.).
   When the root name server directs the query to the correct TLD server, the TLD server then provides information about the authoritative name servers for the requested domain.

3. **Authoritative Name Servers**  
    The **authoritative name servers** hold the final DNS records for the domain.
   They are responsible for returning the IP address of the requested domain.
   These servers can be the domain's hosting provider or a DNS service you configure (such as Cloudflare or Google DNS).
   **This implementation works as an authoritative name server**

### Types of Resolvers

There are two main types of DNS resolvers that interact with these name servers:

- **Recursive Resolver**  
   A **recursive resolver** is responsible for making the full query process.
  It starts at the root name servers and follows the chain through the TLD servers until it
  reaches the authoritative servers. It then returns the final result (IP address) to the client.
  This type of resolver performs all the necessary steps and provides a complete answer.

- **Iterative Resolver**  
   An **iterative resolver**, on the other hand, only provides partial answers.
  It queries the root name servers for a direction to the TLD servers and then asks the TLD servers for the authoritative servers.
  The iterative resolver does not follow the full path and relies on the client (or a recursive resolver) to complete the final step.

## DNS Message Format

The DNS protocol uses a specific binary message format for communication between clients (resolvers) and servers.
Each DNS message consists of a **header**, a **question section**, and optional **answer**, **authority**, and **additional sections**.
Below is a detailed explanation of each part of the DNS message.

---

### DNS Message Structure

1. **Header (12 bytes)**  
   The header contains metadata about the DNS query or response.

   | Field               | Size (bits) | Description                                                  |
   | ------------------- | ----------- | ------------------------------------------------------------ |
   | Transaction ID      | 16          | Unique identifier for matching requests and responses.       |
   | Flags               | 16          | Control and status flags (e.g., QR, Opcode, AA, TC, RD, RA). |
   | Question Count      | 16          | Number of questions in the Question section.                 |
   | Answer Record Count | 16          | Number of resource records in the Answer section.            |
   | Authority Count     | 16          | Number of resource records in the Authority section.         |
   | Additional Count    | 16          | Number of resource records in the Additional section.        |

   #### Header Flags Breakdown (16 bits)

   The flags field is divided into individual control flags and operation codes:
   | Bit | Name | Description |
   |-----|---------------|-----------------------------------------------------------------------------|
   | 0 | QR | Query (0) or Response (1). |
   | 1-4 | Opcode | Type of query (e.g., 0 for standard query). |
   | 5 | AA | Authoritative Answer (set by authoritative servers in responses). |
   | 6 | TC | Truncated Message (set if the message is too large for the transport). |
   | 7 | RD | Recursion Desired (set by clients to request recursion). |
   | 8 | RA | Recursion Available (set by servers that support recursion). |
   | 9-11| Z | Reserved for future use; must be set to 0. |
   | 12-15 | RCODE | Response Code (e.g., 0 for NOERROR, 3 for NXDOMAIN). |

---

2. **Question Section**
   The question section specifies the query details sent by the client.

   | Field  | Size (variable) | Description                                   |
   | ------ | --------------- | --------------------------------------------- |
   | QNAME  | Variable        | Domain name being queried, encoded in labels. |
   | QTYPE  | 16 bits         | Type of query (e.g., A, AAAA, CNAME, MX).     |
   | QCLASS | 16 bits         | Class of query (e.g., IN for Internet).       |

   #### Example:

   For the domain `example.com`, the `QNAME` is encoded as: [7]example[3]com[0]

Where each label is preceded by its length, and the name is terminated by a zero byte.

---

3. **Answer Section**
   Contains resource records (RRs) that answer the query. Each RR has the following fields:

| Field    | Size (variable) | Description                                                |
| -------- | --------------- | ---------------------------------------------------------- |
| NAME     | Variable        | Domain name to which the RR applies (often a pointer).     |
| TYPE     | 16 bits         | Type of RR (e.g., A, AAAA, CNAME).                         |
| CLASS    | 16 bits         | Class of RR (e.g., IN for Internet).                       |
| TTL      | 32 bits         | Time-to-live, in seconds, for caching the RR.              |
| RDLENGTH | 16 bits         | Length of the RDATA field.                                 |
| RDATA    | Variable        | Data associated with the RR (e.g., IP address for A type). |

---

4. **Authority Section**
   Contains RRs pointing to authoritative servers for the queried domain. The format is identical to the Answer section.

---

5. **Additional Section**
   Provides additional data, such as: EDNS cookies or resolved records that might be useful for the resolver

## Name Compression in DNS Messages

To reduce the size of DNS messages, domain names can be encoded using **name compression**. Name compression replaces repeated domain name labels with pointers to their previous occurrences in the message.

### Pointer Format

A pointer uses the most significant 2 bits of a label length byte to indicate a compressed label:

- **Pointer Indicator:** The first two bits of the byte are set to `11` (binary).
- **Offset:** The remaining 14 bits represent the pointer's offset from the start of the DNS message.

### Example:

Consider the domain name `example.com` appearing multiple times in a DNS response:

- The first occurrence is encoded as: [7]example[3]com[0]

- Subsequent occurrences use a pointer: [C0 0C]

## Extended DNS (EDNS)

EDNS extends DNS by adding an **OPT pseudo-record** in the additional section, enabling features like larger message sizes and DNS Cookies.

### OPT Record Structure:

| Field          | Size     | Description                              |
| -------------- | -------- | ---------------------------------------- |
| NAME           | 0 bytes  | Must be 0 (root domain).                 |
| TYPE           | 16 bits  | Always `OPT` (value 41).                 |
| UDP Size       | 16 bits  | Maximum UDP payload size supported.      |
| Extended RCODE | 8 bits   | Extended response codes.                 |
| Version        | 8 bits   | EDNS version (currently 0).              |
| Flags          | 16 bits  | Various flags (e.g., DO bit for DNSSEC). |
| Data           | Variable | Optional data like DNS Cookies.          |

---

## Requirements

- [Go 1.20+](https://go.dev/dl/)
- [PostgreSQL 12+](https://www.postgresql.org/)
- [GORM](https://gorm.io/) ORM for database operations
- [GIN](https://gin-gonic.com/) Web framework for REST API interface
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/XxRoloxX/dns.git
cd dns
```

2. Set up environment variable for database configuration

```bash
export DB_HOST=localhost
export DB_USER=youruser
export DB_PASSWORD=yourpassword
export DB_NAME=yourdb
export DB_PORT=5432
```

or just write them to the .env file

3. Build and run docker-compose

```bash
docker compose up -d
```

![Alt text](./assets/dns-check.gif)

## Deployment

To begin using this server in production, register a NS
(Name Server) record with your domain registrar and point it
to the instance hosting this server. Once this is done,
the server will be able to manage the **delegated zone**, and
all DNS queries for the associated subdomain will be directed to the newly deployed server.

For example, if your domain is `example.com` and the subdomain is `sub.example.com`,
you would add a NS record pointing sub.example.com to the IP address or
hostname of the new server, allowing it to handle all DNS traffic for that subdomain.

## API ENDPOINTS

### BASE URL

`http://localhost:8080/`

### For DNS server

`localhost:53` (udp)

### Endpoints

#### Get all Records

_GET_ `/records`

- Request response

```json
[
  {
    "name": "example.com",
    "type": "A",
    "class": "IN",
    "data": "192.168.1.1"
  }
]
```

#### Create record

_POST_ `/records`

- Request body

```json
[
  {
    "name": "example.com",
    "type": "A",
    "class": "IN",
    "data": "192.168.1.1"
  }
]
```

#### Delete record

_DELETE_ `/records/:id`

```json
{
  "message": "Record deleted succefully"
}
```

![Alt text](./assets/management-check.gif)

## Resources

- [RFC 1035: Domain Names - Implementation and Specification](https://datatracker.ietf.org/doc/html/rfc1035)
- [RFC 1034: Concepts and Facilities](https://www.ietf.org/rfc/rfc1034.txt)
- [RFC 6891: Extension Mechanisms for DNS (EDNS)](https://datatracker.ietf.org/doc/html/rfc6891)
