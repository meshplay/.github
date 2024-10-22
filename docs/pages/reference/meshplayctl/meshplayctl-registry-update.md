---
layout: default
title: meshplayctl-registry-update
permalink: reference/meshplayctl/registry/update
redirect_from: reference/meshplayctl/registry/update/
type: reference
display-title: "false"
language: en
command: registry
subcommand: update
---

# meshplayctl registry update

Update the registry with latest data.

## Synopsis

Updates the component metadata (SVGs, shapes, styles and other) by referring from a Google Spreadsheet.
<pre class='codeblock-pre'>
<div class='codeblock'>
meshplayctl registry update [flags]

</div>
</pre> 

## Examples

Update models from Meshplay Integration Spreadsheet
<pre class='codeblock-pre'>
<div class='codeblock'>
meshplayctl registry update --spreadsheet-id [id] --spreadsheet-cred [base64 encoded spreadsheet credential] -i [path to the directory containing models].

</div>
</pre> 

Updating models in the meshplay/meshplay repo
<pre class='codeblock-pre'>
<div class='codeblock'>
meshplayctl registry update --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --spreadsheet-cred $CRED

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help                      help for update
  -i, --input string              relative or absolute input path to the models directory (default "../server/meshmodel")
      --spreadsheet-cred string   base64 encoded credential to download the spreadsheet
      --spreadsheet-id string     spreadsheet it for the integration spreadsheet

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshplay/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/reference/meshplayctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
