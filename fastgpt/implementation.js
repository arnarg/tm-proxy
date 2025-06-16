async function fetch_fastgpt_content(query, apiKey, pluginServer) {
  const response = await fetch(
    `${pluginServer}/web-search/fastgpt?q=${encodeURIComponent(query)}`,
    {headers: {'Kagi-API-Key': apiKey}}
  );

   if (!response.ok) {
    throw new Error(
      `Failed to fetch search results: ${response.status} - ${response.statusText}`
    );
  }

  const data = await response.json();
  return data.responseObject.content;
}

async function get_fastgpt_search_results(params, userSettings) {
  const { keyword } = params;
  const { pluginServer, kagiAPIKey } = userSettings;

  if (!kagiAPIKey) {
    throw new Error(
      'Please set the API Key in the plugin settings.'
    );
  }

  if (!pluginServer) {
    throw new Error(
      'Missing plugin server URL. Please set it in the plugin settings.'
    );
  }

  const cleanPluginServer = pluginServer.replace(/\/$/, '');

  try {
    return await fetch_fastgpt_content(keyword, kagiAPIKey, cleanPluginServer);
  } catch (error) {
    console.error('Error getting search results:', error);
    return 'Error: Unable to fetch search results. Please try again later.';
  }
}
