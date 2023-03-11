# ccommits_pls

Conventional Commits Language Server (or Conventional Commits "pls")

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
    cmd = {
      "/Users/cristianoliveira/other/ccommits_pls/bin/ccommits_pls",
    },
    filetypes = { "gitcommit" },
    root_dir = util.path.dirname,
    -- init_options = {
    --   command = { "ffff", "fooooo" },
    -- },
    autostart = true,
  },
  -- on_new_config = function(new_config) end;
  -- on_attach = function(client, bufnr) end;
  docs = {
    description = [[
    Language Server Protocol for Conventional Commits.
    ]],
    default_config = {
      root_dir = [[root_pattern(".git")]],
    },
  },
}

local lsp_flags = {
  -- This is the default in Nvim 0.7+
  debounce_text_changes = 150,
}

local servers = {
  dockerls = {},
  gopls = {},
  -- pyright = {},
  rust_analyzer = {},
  tsserver = {},

  lua_ls = {
    Lua = {
      workspace = { checkThirdParty = false },
      telemetry = { enable = false },
      diagnostics = {
        globals = { 'vim' }
      }
    },
  },
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

lspconfig.ccommits_pls.setup {
  on_attach = Lsp_on_attach, -- see ../mappings/lsp.lua
  flags = lsp_flags,
}
``
