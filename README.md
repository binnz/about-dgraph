# about-dgraph

# What Is Dgraph
Dgraph 是一个图数据库，介绍Dgraph之前先来说下简单说下什么是图和图数据库。

在数学的分支图论中，图（Graph）用于表示物件与物件之间的关系，是图论的基本研究对象。一张图由一些小圆点（称为顶点或结点）和连结这些圆点的直线或曲线（称为边）组成（wikipedia）。在知识图谱当中，每个结点用来表示一个实体（人、公司等），边用来表示实体之间的关系，著名的图结构如Google的知识图谱和Facebook的社交图谱。
图数据库首先是一个数据库，它用图结构（包含结点和边）来表示、存储和提供数据。

Dgraph就是一个开源的图数据库。它有着低延迟、高吞吐量的特性，原生支持分布式。

# QuickStart

## Data

```
wget https://github.com/dgraph-io/benchmarks/blob/master/data/1million.schema
wget https://github.com/dgraph-io/benchmarks/blob/master/data/1million.rdf.gz
```

{
  pulpFiction(func:anyofterms(name@en, " Tony")) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
## Deploy

Downloads any version of Dgraph, like v21.3.0 here.
```
https://github.com/dgraph-io/dgraph/releases/download/v21.03.0/dgraph-linux-amd64.tar.gz
```

```
tar -xvf dgraph-linux-amd64.tar.gz
```

├── badger
├── dgraph
├── dgraph-linux-amd64.tar.gz
├── dgraph-ratel

Start zero node.
```
nohup ./dgraph alpha  --zero localhost:5080 -p out/0/p &
```
Bulk load the part of the movie dataset, then generate the p directory.

```
./dgraph bulk -f 1million.rdf.gz -s 1million.schema --map_shards=1 --reduce_shards=1 --http localhost:8000 --zero=localhost:5080
```

Start one alpha node. 
```
nohup ./dgraph alpha  --zero localhost:5080 -p out/0/p &
```

Time to test the queries and mutations.

## Query

### 1. fundamentals

Query Node named "Pulp Fiction"
```
{
  pulpFiction(func: eq(name@en, "Pulp Fiction")) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
```

Query by UID
```
{
  pulpFiction(func: uid(0x5a36b2aa5793b026)) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
```

Multiple IDs can be specified in a list to the uid function
```
{
  pulpFiction(func: uid(0xfffd8d67d8358e15, 0x5a36b2aa5793b026)) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
```

A query expands edges from node to node by nesting query blocks with { }.
if we want to know the actor and character of Pulp Fiction, expand from the starring edge, then the performance.actor and  performance.character edge(we need to understand the dataset first).
```
{
  brCharacters(func: eq(name@en, "Pulp Fiction")) {
    name@en
    initial_release_date
    starring {
      performance.actor {
        name@en  # actor name
      }
      performance.character {
        name@en  # character name
      }
    }
  }
}
```
### 2. Functions we can use
Comparison functions (eq, ge, gt, le, lt) in the query root (aka func:) can only be applied on indexed predicates.
Comparison functions used on @filter directives even on predicates that have not been indexed.
All other functions, in the query root or in the filter can only be applied to indexed predicates.

*** 
```
matches nodes with an outgoing string field fieldName where the string contains all listed terms.


Syntax: allofterms(predicate, "space-separated term list")
```

Used in root: Query All nodes that have either Science or Fiction in the name.
```
{
  pulpFiction(func: allofterms(name@en, "kill nill")) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
```
Used in filter:Query all movies of Quentin Tarantino that has kill and bill in film name.
@filter(has(director.film)) means return only director(has director.film edge)
```
{
  me(func: eq(name@en, "Quentin Tarantino"))@filter(has(director.film)) {
    name@en
    director.film @filter(allofterms(name@en, "kill bill"))  {
      name@en
    }
  }
}
```
***
```
matches nodes with an outgoing string field fieldName where the string contains at least one of the listed terms.

Syntax Example: anyofterms(predicate, "space-separated term list")
```

Usage at root: any node has either Science or Fiction in name edge.
```
{
  pulpFiction(func: anyofterms(name@en, "Science Fiction")) {
    uid
    name@en
    initial_release_date
    netflix_id
  }
}
```

Usage at filter: query all node has edge director.film and name edge value is Quentin Tarantino. Return only director.film edges has brown or bill in name edge.
```
{
  me(func: eq(name@en, "Quentin Tarantino"))@filter(has(director.film)) {
    name@en
    director.film @filter(anyofterms(name@en, "brown bill"))  {
      name@en
    }
  }
}
```
***
```
Fuzzy matching. Matches predicate values by calculating the Levenshtein distance to the string

Syntax: match(predicate, string, distance)
```
```
{
    bill(func:match(name@en, bill, 1)){
	name@en
    }
}
```
***
```
Full-Text Search. Apply full-text search with stemming and stop words to find strings matching all or any of the given text.

The following steps are applied during index generation and to process full-text search arguments:

    Tokenization (according to Unicode word boundaries).
    Conversion to lowercase.
    Unicode-normalization (to Normalization Form KC).
    Stemming using language-specific stemmer (if supported by language).
    Stop words removal (if supported by language).


Syntax: alloftext(predicate, "space-separated text")
        anyoftext(predicate, "space-separated text")
```
```
{
  movie(func:alloftext(name@en, "the dog which barks")) {
    name@en
  }
}

{
  movie(func:anyoftext(name@en, "the dog which barks")) {
    name@en
  }
}
```
***
```
Equal to

Syntax:

eq(predicate, value)
eq(val(varName), value)
eq(predicate, val(varName))
eq(count(predicate), value)
eq(predicate, [val1, val2, ..., valN])
eq(predicate, [$var1, "value", ..., $varN])
```
```
{
  steve as var(func: anyofterms(name@en, "bill")) {
    films as count(director.film)
  }

  stevens(func: uid(steve)) @filter(eq(val(films), [1,2,3])) {
    name@en
    numFilms : val(films)
  }
}
```
here about query blocks and value Variables will detailed later.
***
```
Less than, less than or equal to, greater than and greater than or equal to
Syntax : for inequality IE

IE(predicate, value)
IE(val(varName), value)
IE(predicate, val(varName))
IE(count(predicate), value)
With IE replaced by

le less than or equal to
lt less than
ge greater than or equal to
gt greater than
```
***
```
between
Syntax: between(predicate, startDateValue, endDateValue)
```
```
{
  me(func: between(initial_release_date, "1989-01-01", "2007-12-31")) @filter(anyofterms(name@en,"bill")){
    name@en
    genre {
      name@en
    }
  }
}
```
***
```
uid
Syntax:
    q(func: uid(<uid>))
    predicate @filter(uid(<uid1>, ..., <uidn>))
    predicate @filter(uid(a)) for variable a
    q(func: uid(a,b)) for variables a and b
    q(func: uid([]))
```
```
{
	me(func:uid(0x6e3c7860e962104c, 0xbaababda58934246)){

  name@en
  }
  }


{
	me(func:anyofterms(name@en, "Lee")){
	name@en
    director.film @filter(uid(0x6e3c7860e962104c, 0xbaababda58934246)){
    name@en
    }
    }
}
```

# GraphQL
GraphQL is a data query language developed internally by Facebook in 2012 before being publicly released in 2015. It provides an alternative to RESTful architectures.(wikipedia)
先说下RESTFul规范。RESTful是web开发当中客户端和服务端数据交互的一种规范，通过定义一系列的约束来实现；满足这种规范的架构设计称为RESTful风格的架构，满足这种规范的web服务被称为RESTful web服务。
Graph是一种查询语言或语法规范，也是用来描述客户端如何向服务端请求数据。与RESfFul类似，它的特点是客户端可以通过一种语言准确描述所需要的数据，方便用一个请求获取到多个数据源的聚合数据，而不用发送多个RESTful请求。
# DQL

filter

Dgraph has many functions available for filtering results, including the following:

fieldName: { allofterms: "term1 ... termN" }: matches nodes with an outgoing string field fieldName where the string contains all listed terms

fieldName: { anyofterms: "term1 ... termN" }: As with allOfTerms, but matches at least one term

The equalities and inequalities can be applied to edges of types: Int, Int64, Float, String and DateTime, including:

fieldName: { eq: value }: equal to
fieldName: { ge: value }: greater than or equal to
fieldName: { le: value }: less than or equal to
fieldName: { gt: value }: greater than
fieldName: { lt: value }: less than





{
  getPerson(xid: "alice") {
    xid
    friends(filter: {
      and: [
        { age: { ge: 30 } },
        { age: { lt: 40 } }
      ]
    }) {
      xid
      age
    }
  }
}


{
  getPerson(xid: "alice") {
    xid
    friends(filter: {
      age: { ge: 35 }, or: { xid: { eq: "frank" } }
    }) {
      xid
      age
    }
  }
}



The not logical wrapper cannot contain a list of arguments not: [{ age: ... }, ...]. But it can contain multiple arguments and other logical combinators, e.g.: not: { and: [{ age: ... }, ...] }


Sort
{
  getPerson(xid: "alice") {
    id
    name
    friends(order: { desc: age, then: { asc: name } }) {
      id
      name
      age
    }
  }
}




Dgraph generates first and offset arguments that can be used in combination to achieve such limits and paginate results:

first: N Return only the first N results
offset: N Skip the first N results
By default, query answers are ordered by uid. You can change this behavior by explicitly specifying an order.

The first and offset arguments are available on query<Type> queries and on any edge to a list of nodes.



aggregate
count - count how many friends Alice has
xidMin - find the minimum xid value sorted alphabetically
xidMax - find the maximum xid value sorted alphabetically
nameMin - find the minimum name value sorted alphabetically
nameMax - find the maximum name value sorted alphabetically
ageMin - find the minimum age of Alice’s friends
ageMax - find the maximum age of Alice’s friends
ageAvg - find the average age of Alice’s friends
ageSum - sum of all of the ages of Alice’s friends


How Dgraph Search Works
The graphs in Dgraph can be huge, so starting your search from all nodes isn’t always efficient. Dgraph needs a place to start searching, and you choose this by starting your search from a root query function.

The function has(edge_name) returns nodes that have an outgoing edge of the given name.

{
  queryPerson(filter: { has: friends, and: { has: age } }) {
    id
    name
    age
    friendsAggregate(filter: { has: age }) {
      count
    }
  }
}

{
  users: queryPerson {
    username: xid
  }
  peopleWithPets: queryPerson(filter: { has: ownsPets }) {
    name
    pet: ownsPets {
      name
    }
  }
}

cascade:
{
  getPerson(xid: "alice") @cascade {
    name
    age
    friends {
      name
      ownsPets {
        name
      }
    }
  }
}

{
  getPerson(xid: "alice") {
    name
    age
    friends @cascade {
      name
      ownsPets(filter: { xid: { eq: "goldie" } }) {
        name
      }
    }
  }
}
# Slash GraphQL
# Clients
# Deploy


# References

[graphql-fundamentals](https://dgraph.io/docs/query-language/graphql-fundamentals/)
# 幕后故事（作者、本项目历程）


	
<img src="/pic/anyofterm.png" width="200" height="200"/><br/>