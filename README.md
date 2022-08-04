# mw2md
## A dead simple tool for converting a MediaWiki backup to Markdown files
Useful for moving from MediaWiki to other wiki platforms which use Markdown
Requires pandoc to be installed in PATH

Usage:

Get a export from Special:Export and run:

```
mw2md [YOUR EXPORT]
```
## Notes
If a page has subpages, it will currently convert the '\' to -. A directory solution might be worked in the future.

It doesn't currently support revisions.

# License
Licensed under the MIT license.