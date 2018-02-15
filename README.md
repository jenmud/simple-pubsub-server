PubSub Server
=============

This is a quick mockup of a type of pub/sub server. It is a very very simple implementation of a basic publish and subscribe service.

**NOTE: This is not production tested, and it is not performance tested!**
**It was a learning project that I started in my spare time on the train.**

Expected behaviour
==================

Publisher
---------

Once a client is subscribed as a publisher to a topic, you can **not** send other
messages such as *ping*, *subscribe*, *disconnect*. The reason for this is because
we are simply fowarding on all messages sent by the publishing client and never get
the change to be dispatched properly. This is a unfortunate outcome from a poor
design.

There can only ever be a single client publishing to a topic.


Subscriber/Consumer
-------------------

A subscriber/consumer client can subscribe to as many topic as it wants. It can only subscribe to a topic that is being published.

Clients will need to re-subscribe to topics if the publishers have disconnected and
reconnected.

You can have any amount of subscribers/consumers to a single topic as you want.


Heartbeats
----------

Each client once connected will periodically receive a *heartbeat* message. A *heartbeat* message will have a *Tick* attribute which is the time when the message
was sent. You can use this to workout if the server is still working as expected.


Example of publishing
=====================

With the server running and accepting connections

```bash
$ ./myPubSub
2018/02/15 21:15:39 Server running and accepting connections on port 3003
```

Connect to the server and issue the *connect* message and you will receive a *Welcome* message

```bash
$ telnet localhost 3003
Trying ::1...
Connected to localhost.
Escape character is '^]'.
<connect />
<Welcome name="" address="[::1]:51713" datetime="2018-02-15 21:16:24.576339 +1100 AEDT m=+45.094943292"><topics></topics></Welcome>
```

Now send the *publish* message with a topic name

```bash
<publish topic="demo" />
```


Example of subscribing/consuming
--------------------------------

With a client already publishing to a topic (see above), connect another client to the server

```bash
$ telnet localhost 3003
Trying ::1...
Connected to localhost.
Escape character is '^]'.
<connect />
<Welcome name="" address="[::1]:51751" datetime="2018-02-15 21:23:18.273907 +1100 AEDT m=+458.780100610"><topics><topic>demo</topic></topics></Welcome>
```

```bash
<subscribe topic="demo" />
```


Supported messages
==================

| Message        | Description                                      | Required |
| -------------  |:------------------------------------------------:| --------:|
| connect        | Connect and register the client with the server. | yes      |
| ping           | Send a ping and receive a pong reply.            | no       |
| subscribe      | Subscribe to a topic.                            | no       |
| publish        | Register the client as a publisher of a topic.   | no       |
| disconnect     | Disconnect and unregister client from the server | no       |