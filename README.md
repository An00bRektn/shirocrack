# shirocrack
> Simple hash cracker for [Apache Shiro](https://github.com/apache/shiro) hashes written in Golang. Useful for exploiting [CVE-2024-4956](https://github.com/erickfernandox/CVE-2024-4956).

## Let's keep this short

- CVE-2024-4956 is a path traversal vulnerability in Sonatype Nexus Repository that allows an unauthenticated attacker with network access to the server to read arbitrary files from the system, as long as they know the path
- Every exploit demo I saw pulled `/etc/passwd` and called it a day, which is good, but doesn't *really* explain the critical CVSS score it got
    - Nexus 3.60.0, in its default setup, uses OrientDB, a NoSQL database that I honestly don't fully understand yet
        - One part of OrientDB is the idea of "clustering", which for our purposes, means they store data in these binary `.pcl` files
    - Nexus stores `.pcl` files in the following directories
        - `/nexus-data/db/OSystem`
        - `/nexus-data/db/component`
        - `/nexus-data/db/config`
        - `/nexus-data/db/security`
    - `/nexus-data/db/security/user*.pcl` contains password hashes
- And now for why this exists: Turns out hashcat doesn't support the format, and the only cracker I found online was written in Java (bad): [GitHub Gist - gquere](https://gist.github.com/gquere/365cfcceef9ac8d145cc59bbf2c27648)
    - No shade to the creator, I just needed something faster
- I rewrote that in Golang, which should be good enough.

## Usage
```shell
$ ./shirocrack
[i] usage: ./shirocrack HASH wordlist.txt

$ ./shirocrack '$shiro1$SHA-512$1024$+rU2PizvJ/Nj7s4XDn866A==$5fGRXQstvAgoVA1N8ipEYzsQFFN8VqmNLsKs/Ka8x1FrxflDaxXprx/vwLhZDBOXABT72E0H/SNpnQSLQgW87g==' ./wordlist.txt
[+] Found match: 5fGRXQstvAgoVA1N8ipEYzsQFFN8VqmNLsKs/Ka8x1FrxflDaxXprx/vwLhZDBOXABT72E0H/SNpnQSLQgW87g==:an00brektn
[+] Success!
```
