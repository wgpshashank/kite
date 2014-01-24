New Kite Framework
==================

This document tries to describe the new Kite framework and components involved.

What is a Kite?
---------------

A Kite is a **micro-service** talking a special *Kite Protocol*. Installed
browser apps on koding.com talks with these Kites. For example, file tree in
*Develop* tab talks with *fs* Kite.

What can a Kite do?
-------------------

A Kite can respond to requests and make requests to other Kites. The request
may return a response or it may call a callback that is sent in the arguments
of the request any time. The details are explained in *Kite Protocol*.

What is Kite Protocol?
----------------------

Kite Protocol is the subset of **dnode** protocol. Information about the dnode
protocol can be found [here](https://github.com/substack/dnode-
protocol/blob/master/doc/protocol.markdown). It is a JSON based protocol which
allows us to build async-RPC and pub/sub mechanisms on top of it. Read and
understand it first. Here are the differences in our protocol:

* We do not use the initial "methods" exchange.
* We do not use the "links" field in the message.
* We always send 1 or 2 arguments in "arguments" field:
    * The first one is *options* object, it is an object with 3 keys:
        * **kite:** Information about the Kite that is making the request.
        * **authentication:** Authentication information for this request.
        * **withArgs:** This fields contains the arguments that will be passed to the method called. It may be a type of any valid JSON type.
    * The second one is the optional response callback function. It may be omitted or passed as *null*. We remove the callback after we get the response to free memory because it will only be called once by the other side.
* Currently, we use WebSocket for transport layer. This may change in future. The path to connect on the server is "/dnode".


What does a Kite Protocol message looks like?
---------------------------------------------

This is a message for geting the list of kites from *Kontrol*:

(Contents of `withArgs` are specific to `getKites` method and can be totally different for any other method.)

```js
{
  "arguments":[
    {
      "authentication":{
        "key":"KpVmfxow4IJCQ3GiRXi3PWu2OMitnDhrLWTrm0TYMqSRwvFL0GQsKiouBL889Iu9",
        "type":"kodingKey"
      },
      "kite":{
        "environment":"development",
        "hostname":"tardis.local",
        "id":"d1a89409-079d-43fc-44d0-999a111c0915",
        "name":"application",
        "port":"52649",
        "publicIP":"",
        "region":"localhost",
        "username":"devrim",
        "version":"1"
      },
      "withArgs":{
        "environment":"",
        "hostname":"",
        "id":"",
        "name":"mathworker",
        "region":"",
        "username":"devrim",
        "version":""
      }
    },
    "[Function]"
  ],
  "callbacks":{
    "1":[
      "1"
    ]
  },
  "links":[

  ],
  "method":"getKites"
}
```

...and this is the response to the request above:

```js
{
  "arguments":[
    null,
    [
      {
        "kite":{
          "environment":"development",
          "hostname":"tardis.local",
          "id":"80663a74-2715-4cc4-5392-36d5cbdb0434",
          "name":"mathworker",
          "port":"52641",
          "publicIP":"127.0.0.1",
          "region":"localhost",
          "username":"devrim",
          "version":"1"
        },
        "token":"IQhRVGIaclwBjAs7FujrbsJeSGwUSOqsoOS8u-f0uFoSyvmjvYHG4wgpdcS36sxixP1giH3nn--JDzA4KMuupZnkyZzi0NbCb_ktX-asGWPzyhU45lvPuBRxVvgQKHTr6uD9kmxhgmbdN3sohcX7JpqTlkuJSVnGwLyFio6s4NkAYcQ="
      }
    ]
  ],
  "callbacks":{

  },
  "links":[

  ],
  "method":1
}
```

Implementation of sending a message: [go/src/koding/newkite/kite/remote.go](https://git.portal.sj.koding.com/koding/koding/blob/master/go/src/koding/newkite/kite/remote.go)

Implementation of processing incoming message: [go/src/koding/newkite/kite/request.go](https://git.portal.sj.koding.com/koding/koding/blob/master/go/src/koding/newkite/kite/request.go)

What is *Kontrol*?
------------------

Kontrol is a special Kite that is run by Koding. It is both a dynamic name service (like DNS) and an authentication service (like OAuth2). We are thinking about seperating these services into 2 different Kites in future.

Kontrol code is here: [go/src/koding/newkite/kontrol/main.go](https://git.portal.sj.koding.com/koding/koding/blob/master/go/src/koding/newkite/kontrol/main.go)

How a Kite can find other Kites?
--------------------------------

A Kite needs the address of a Kite to connect to it. The address can be asked to *Kontrol* with a query. It's like DNS. Kontrol is also a Kite that has following methods:

* **register:** A kite need to call this method after it has connected to Kontrol in order other Kites to find it. The method takes no arguments and returns error as a response.

* **getKites:** A kite issues a query to get the information about the other Kites including the addresses.

* **watchKites:** Takes the same query parameter as *getKites* but does not return the list of queried kites. Instead, returns a single Kite as they register to *Kontrol*.

What does a *query* look like?
----------------------------

```go
type KontrolQuery struct {
    Username    string `json:"username"`
    Environment string `json:"environment"`
    Name        string `json:"name"`
    Version     string `json:"version"`
    Region      string `json:"region"`
    Hostname    string `json:"hostname"`
    ID          string `json:"id"`
}
```

* The order of the fields is from general to specific.
* Missing fields mathes everything.


How is authentication handled?
------------------------------

Every request must have an authentication information sent in *options.authentication* field in the message. There are different authentication types:

* **kodingKey:** Used to talk with *Kontrol*. This is a secret key that is unique to users machine and created by **kd tool**. Since both sides know the key *Kontrol* can authenticate the user from this key.

* **sessionID:** Browser clients must use this authentication type to talk with *Kontrol*.

* **token:** Used in *Kite-to-Kite* communication. When a Kite requests the list of other Kites from *Kontrol* by invoking it's `getKites` method, *Kontrol* also includes auto-generated, self-expiring tokens in the response. These tokens are valid for only one Kite and can be validated only on the Kite that is the receiver of the request.

What are the basic services that every Kite must offer?
------------------------------------------------------------

Currently all Kites must support these 2 methods below. You don't have to
implement them, they are implemented in Kite framework.

* **heartbeat:** A heartbeat service that calls the given callback function on
**given interval. If it cannot make the call, *Kontrol* unregisters the
**Kite.

* **vm.info:** To get information about the running host (cpu, disk, etc.)

What about the browser?
-----------------------

Browser is also acts like a Kite. It has it's own methods ("log" for logging a
message to the console, "alert" for displaying alert to the user, etc.). A
connected Kite can call methods of the browser. See how it's implemented in [c
lient/app/MainApp/kite/newkite.coffee](https://git.portal.sj.koding.com/koding
/koding/blob/master/client/app/MainApp/kite/newkite.coffee).

*Kontrol client can be inspected to see how `NewKite` class can be used in [cl
*ient/app/MainApp/kite/kontrol.coffee](https://git.portal.sj.koding.com/koding
*/koding/blob/master/client/app/MainApp/kite/kontrol.coffee).

How can I write a new Kite?
---------------------------

* Import the Kite Framework.
* Add your method handlers.
* Call `Kite.Run()`

See [go/src/koding/newkite/examples/mathworker.go](https://git.portal.sj.kodin
g.com/koding/koding/blob/master/go/src/koding/newkite/examples/mathworker.go)
for an example Kite code.

Or read some real Kite code in [go/src/koding/kites](https://git.portal.sj.kod
ing.com/koding/koding/tree/master/go/src/koding/kites).