# tm-proxy

An alternative typingmind plugin proxy.

## Why?

- I had some problems with packaging the official one for NixOS.
- The `web-page-reader` endpoint was stripping out hyperlinks in the fetched web page.
- The endpoint is well known and it's too easy for anyone to spam the `web-page-reader` endpoint on my self-hosted instance if exposed publicly (which I needed to do).
  - I solved this be specifying a URI prefix (i.e. `/<some-random-string>/web-page-reader/get-content`) so the endpoint is less discoverable (security through obscurity).

## Using

### Web Page Reader

The plugin server works with the standard `Web Page Reader` plugin. Just set the `Plugin Server` URL to your instance with `https://<domain>/[<prefix>]`.

### FastGPT Search

Create a new plugin, using the JSON editor paste the following:

```json
{
    "id": "get_fastgpt_search_results",
    "code": "async function fetch_fastgpt_content(query, apiKey, pluginServer) {\n  const response = await fetch(\n    `${pluginServer}/web-search/fastgpt?q=${encodeURIComponent(query)}`,\n    {headers: {'Kagi-API-Key': apiKey}}\n  );\n\n   if (!response.ok) {\n    throw new Error(\n      `Failed to fetch search results: ${response.status} - ${response.statusText}`\n    );\n  }\n\n  const data = await response.json();\n  return data.responseObject.content;\n}\n\nasync function get_fastgpt_search_results(params, userSettings) {\n  const { keyword } = params;\n  const { pluginServer, kagiAPIKey } = userSettings;\n\n  if (!kagiAPIKey) {\n    throw new Error(\n      'Please set the API Key in the plugin settings.'\n    );\n  }\n\n  if (!pluginServer) {\n    throw new Error(\n      'Missing plugin server URL. Please set it in the plugin settings.'\n    );\n  }\n\n  const cleanPluginServer = pluginServer.replace(/\\/$/, '');\n\n  try {\n    return await fetch_fastgpt_content(keyword, kagiAPIKey, cleanPluginServer);\n  } catch (error) {\n    console.error('Error getting search results:', error);\n    return 'Error: Unable to fetch search results. Please try again later.';\n  }\n}",
    "uuid": "24a4a0c9-6b47-4cc4-8016-20f4f1a51fd2",
    "title": "FastGPT Search",
    "iconURL": "https://help.kagi.com/assets/kagi-logo.Bh8O11VU.png",
    "githubURL": "https://github.com/arnarg/tm-proxy",
    "openaiSpec": {
        "name": "get_fastgpt_search_results",
        "parameters": {
            "type": "object",
            "required": [
                "keyword"
            ],
            "properties": {
                "keyword": {
                    "type": "string",
                    "description": "The search keyword"
                }
            }
        },
        "description": "Search for information from the internet in real-time using Kagi FastGPT."
    },
    "outputType": "respond_to_ai",
    "userSettings": [
        {
            "name": "pluginServer",
            "label": "Plugin Server",
            "required": true,
            "description": "The URL of the plugin server",
            "placeholder": "https://..."
        },
        {
            "name": "kagiAPIKey",
            "type": "password",
            "label": "Kagi API Key",
            "required": true
        }
    ],
    "overviewMarkdown": "This plugin allows the AI assistant to search for information from the internet in real-time using Kagi FastGPT.\n\n**🔑 Kagi API Key needed**.\n\nExample usage:\n\n> What's the gold price?\n\n> How's the weather at HCMC at the moment?\n",
    "authenticationType": "AUTH_TYPE_NONE",
    "implementationType": "javascript",
}
```

In the settings set the `Plugin Server` URL to your instance with `https://<domain>/[<prefix>]` and Kagi API Key to an API key acquired from [Kagi](https://kagi.com/settings?p=api).
