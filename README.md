# Moldova
Moldova is a lightweight template interpreter, used to generate random values that plug into the template, as defined by a series of custom tokens.

It understands the tokens as defined further down in the document.

Moldova also comes with a binary executable you can build, which lets you pipe to STDOUT a stream of templates being rendered, broken by a newline. In this way, you could pipe the output to something like Slammer, outlined below.

Moldova comes both as a library you can import, as well as a binary executable that you
can run, which will output results to STDOUT.

## Works great with the Slammer

Moldova was originally designed as input to a database load testing tool, called the [Slammer](https://github.com/StabbyCutyou/slammer). It's primary purpose was to
generate INSERT and SELECT statements, which were then loaded into the Slammer. If you're looking for a good tool to
test the throughput and latency of database operations at scale, checkout the Slammer.

And if you're interested in using them together, checkout the [Moldovan Slammer](http://github.com/StabbyCutyou/moldovan_slammer), a quick helper repo to
demonstrate using them together

## Moldova command

To use the command, first install

```bash
go install github.com/StabbyCutyou/moldova/cmd/moldova
```

The command accepts 2 arguments:

* n - How many templates to render to STDOUT. The default is 1, and it cannot be less than 1.
* t - The template to render

## Example

```bash
moldova -t "INSERT INTO floof VALUES ('{guid}','{guid:ordinal:0}','{country}',{int:min:-2000|max:0},{int:min:100|max:1000},{float:min:-1000.0|max:-540.0},{int:min:1|max:40},'{now}','{now:ordinal:0}','{country:case:up}',NULL,-3)" -n 100
```

This would provide sample output like the following:

```sql
...
INSERT INTO floof VALUES ('791add99-43df-44c8-8251-6f7af7a014df','791add99-43df-44c8-8251-6f7af7a014df','MU',-1540,392,-624.529332,39,'2016-01-24 23:42:49','2016-01-24 23:42:49','UN',NULL,-3)
INSERT INTO floof VALUES ('0ab4cc33-6689-404f-a801-4fd431ca3f30','0ab4cc33-6689-404f-a801-4fd431ca3f30','PL',-1707,112,-550.333145,1,'2016-01-24 23:42:49','2016-01-24 23:42:49','SS',NULL,-3)
INSERT INTO floof VALUES ('a3f4151a-a304-4190-a3df-7fd97ce58588','a3f4151a-a304-4190-a3df-7fd97ce58588','CM',-1755,569,-961.122173,25,'2016-01-24 23:42:49','2016-01-24 23:42:49','NE',NULL,-3)
```

# Tokens

Tokens are represented by special values placed inside of { } characters.

Each token follows this pattern:

{command:argument_name:value|argument_name:value|...|...}

Each token defines it's own set of arguments, which are outlined below. If you're provided multiple arguments, the key-value pairs of argument names to values are delimited by a | pipe character. These can be provided in any order, so long as they are valid arguments.

It is only necessary to provide a | when you have further arguments to provide. For example:

* {integer}
* {integer:min:10}
* {integer:min:10|max:50}
* {integer:min:10|max:50|ordinal:0}


## {guid}

### Options
* ordinal : integer >= 0

### Description

Slammer will replace any instance of {guid} with a GUID/UUID

If you provide the *ordinal:* option, for the current line of text being generated,
you can have the Slammer insert an existing value, rather than a new one. For
example:

"{guid} - {guid:ordinal:0}"

In this example, both guids will be replaced with the same value. This is a way
to back-reference existing generated values, for when you need something repeated.

## {now}

### Options
* ordinal : integer >= 0

### Description

Slammer will replace any instance of {now} with a string representation of Golangs
time.Now() function, formatted per the golang date format example: "2006-01-02 15:04:05".

{now} also supports the *ordinal:* option

## {integer}

### Options
* min : integer < max
* max : integer > min
* ordinal : integer >= 0

### Description

Slammer will replace any instance of {integer} with a random int value, optionally between the range provided. The defaults, if not provided, are 0 to 100.

{integer} also supports *ordinal:* option

## {float}

### Options
* min : integer < max
* max : integer > min
* ordinal : integer >= 0

### Description

Slammer will replace any instance of {float} with a random Float64, optionally between the range provided. The defaults, if not provided, are 0.0 to 100.0

{float} also supports *ordinal:* option

## {unicode}

### Options
* length : integer >= 1
* case : "up" or "down"
* ordinal : integer >= 0

### Description

Slammer will replace any instance of {unicode} with a randomly generated set of unicode
characters, of a length specified by :number. The default value is 2.

{unicode} also takes the :case argument, which is either 'up' or 'down', like so

{unicode:case:up}
{unicode:case:down}

{unicode} also supports *ordinal:* option

Only a certain subset of unicode character ranges are supported by default, as defined
in the moldova/data/unicode.go file.

## {country}

### Options
* case : "up" or "down"
* ordinal : integer >= 0

### Description

Slammer will replace any instance of {country} with an ISO 3166-1 alpha-2 country code.

{country} supports the same *case:* argument as {unicode}. The default value is "up"

{country} also supports the *ordinal:* argument.

# Roadmap

I'll continue to add support for more random value categories, such as a general {time} field, as well as additions to existing ones (for example, a timezone param for :now, as well as the ability to choose a formatting method).

Add proper support for escaping control characters, so they are not interpreted as part of a token. These characters are {, }, :, and |

The current approach is also a bit slow, parsing and interpreting the template each time it's invoked. I'm going to re-write this part of it, so that it builds up an internal callstack of functions during the initial parse, then iterates and invokes the callstack in order to generate N iterations of the template itself, preventing unnecessary re-parses of the template.

# License

Apache v2 - See LICENSE
