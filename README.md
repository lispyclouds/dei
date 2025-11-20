# dei!

This is a (growing and evolving) collection of useful CLI tools that I used pretty much on daily and even hourly basis and has been put together in a nicely packaged and reusable manner as some of the stuff seems to be useful to you all too!

## The name

Dei, or the informal way to call someone with a "Hey" (more [info](https://www.quora.com/What-does-the-Tamil-word-dei-mean-What-are-the-other-such-words-which-we-can-understand-only-by-slang-What-is-an-easy-way-to-learn-to-speak-Tamil)) is something that was common where I come from and it happens to also align with my last name "De".
So whenever someone wanted something, `dei/de` became the de-facto standard to ask me, specially by close friends, something that I truly cherish.
Since I essentially wanted "myself in code" its only natural to call this tool `dei`.

## Usage

### Installation

To install as a Go binary:

- Ensure Go 1.25+ is [installed](https://go.dev/doc/install)
- Run `GOEXPERIMENT=jsonv2 go install github.com/lispyclouds/dei@latest`
- See what's possible: `dei --help`

TODO: Release versioned binaries

### Tools

#### Stateless passwords

In my (not so humble) opinion, passwords should:

- Not need me to remember more than one main password
- Be managed by a tool which:
  - Is open source
  - Simple
  - Stateless
  - Uses peer reviewed and strong cryptographic primitives
  - Customisable
  - Have fast and simple UX

[Spectre](https://spectre.app/) checks all of these boxes for me, specially from the cryptographic point of view. It however doesn't quite have the greatest CLI flow so here we are.

dei implements v3 of the [algorithm](https://spectre.app/spectre-algorithm.pdf) along with an intuitive UX that I think is useful to all users. It implements aggressive caching to speedup the whole process and is optimised to be simple and nimble. Run `dei pw --help` to see all the options.

**Caveat: By default, dei saves the intermeidate scrypt digest to the db as a caching mechanism which also means its there on your local disk UNSECURED. Approaches like PIN encryption, FIDO/hardware token auth etc are being explored for this to mitigate it. Pass --no-cache to opt out and keep the session ephemeral.**

Generate a password for a site, eg: github.com

```bash
dei pw --full-name "Your Full Name" --site "github.com" # pass --to-clipboard to directly copy to clipboard
```

It's recommended to pass the base host name to --site without the protocols, www, paths etc. dei will do a best effort extraction of the host from what is passed but will pass through if it can't.

As mentioned in the algorithm paper, spectre takes in optional parameters like --counter, --class and --variant.
All these options have a default value and when dei encounters a new site, it saves these values to an internal DB and used from there subsequently.
If there is a need to update them, eg --counter 2, pass these explicitly on the CLI and dei with notice the diff and update.

```bash
dei pw --full-name "Full Name" --site "github.com" --counter 3 --class long # update the saved counter and password class
```

#### Crafting conventional commits

[Conventional Commits](https://www.conventionalcommits.org/) are not only a great idea, but is sometimes a mandate in varios projects. Although there are various tools that implement this, I didn't quite like their UX, specially when it comes to caching the right bits to speed not only the tool but your workflow too.

**This is limited to only Git for now.**

dei just helps in two ways:

- help craft the right conventional commit
- help manage a set of co-authors and help attribute them correctly in the message

It does not go any other git interaction like staging, pull, push etc.

Make a commit prompting for the data and co-authors(if configured):

```bash
dei commit
```

Manage co-authors

```bash
dei commit co-authors add --name foo --email bar
dei commit co-authors remove --email bar
dei commit co-authors list
```

## License

Copyright © 2025- Rahul De

Distributed under the MIT License. See LICENSE.
