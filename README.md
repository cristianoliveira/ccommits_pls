# ccommits_pls [![CI Checks](https://github.com/cristianoliveira/ccommits_pls/actions/workflows/checks.yml/badge.svg)](https://github.com/cristianoliveira/ccommits_pls/actions/workflows/checks.yml)

Conventional Commits Language Server (or Conventional Commits "pls") (inspired by [gopls](https://github.com/golang/tools/tree/master/gopls))

## Running for testing

Make sure you have golang version 1.18 

```bash
go get

make build

echo $PWD/bin/ccommits_pls
```

Copy the result and in your editor configure a new lsp Eg. `ccommits_pls`. 

### Neovim

Using mason is super simples

```lua
local lspconfig = require('lspconfig')
local mason_lspconfig = require("mason-lspconfig")

local util = require("lspconfig.util")
local configs = require("lspconfig.configs")

configs.ccommits_pls = {
  default_config = {
    -- Paste here the path to the lsp bin
    cmd = {
      "/Users/youruser/ccommits_pls/bin/ccommits_pls",
    },
    filetypes = { "gitcommit" },
    root_dir = util.path.dirname,
    autostart = true,
  },
  docs = {
    description = [[
    Language Server Protocol for Conventional Commits.
    ]],
    default_config = {
      root_dir = [[root_pattern(".git")]],
    },
  },
}

lspconfig.ccommits_pls.setup {
  on_attach = Lsp_on_attach, -- see ../mappings/lsp.lua
  flags = lsp_flags,
}


-- This section depends on your configuration.
local servers = {
  gopls = {},
}

mason_lspconfig.setup {
  ensure_installed = vim.tbl_keys(servers)
}

mason_lspconfig.setup_handlers {
  function(server_name)
    lspconfig[server_name].setup {
      on_attach = Lsp_on_attach, -- see ../mappings/lsp.lua
      settings = servers[server_name],
      flags = lsp_flags,
    }
  end,
}
``
