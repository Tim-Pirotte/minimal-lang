# Minimal
The compiler for Minimal

## Project Structure
- **cmd:** Application entrypoints
- **built-in:** Implements the logic of the plugins
  - **messenger:** Sends messages that should be shown to the user via output channels
    - **outputs:** The different channels to show the messages
      - **log-renderer:** Render messages in the terminal
      - **test-output:** Captures messages for testing
  - **lexer:** Converts text into a series of tokens
    - **matchers:** Packages that match on textual input
      - **identifier:** Matches text that doesn't start with a digit
      - **indentation:** Keeps track of indentation, newlines and blocks
      - **number:** Matches sequences starting with a number
      - **string:** Matches strings with interpolation and multiline support
      - **raw-string:** Like string but takes the content as is
      - **symbol:** Matches exact strings with a trie for keywords and symbols
      - **space:** Ignores white space
      - **comment:** Matches all text until a new line starting with a prefix
  - **ast:** A pre-order linearized tree with support for diffing and traversing the AST
  - **parsers:** A collection of packages to turn the token stream into an AST
    - **prefix:** Matches on certain token prefixes and runs the attached parsers
    - **prattparser:** Parses with binding power
      - **prefix-unary:** Parses unary expressions with a certain prefix symbol
      - **postfix-unary:** Parses unary expressions with a certain postfix symbol
      - **binary:** Parses binary expressions with a symbol in between
      - **eol:** Wrapper for the prattparser to allow multi-line expressions if the symbol is not ambiguous
        - **grouping:** Allows parenthesizing expressions for precendence and multi-line expressions
  - **symbol-table:** Tracks symbols with interned strings
  - **type-checker:** Passes typing info through the AST in both directions
  - **source:** Contains compilation sources so context in messages can be constructed
  - **substring:** Checks if a string is address-wise a subset of another string
  - **pipeline:** Parallelizes steps within the same compilation stage
  - **ansi:** Type safe way to get ANSI colors
  - **diff:** Generic functions to diff arbitrary data
  - **plugins:** Registry to keep track of all plugins and allow plugins to look for eachother


*The docs of the language itself can be found at: TODO*

