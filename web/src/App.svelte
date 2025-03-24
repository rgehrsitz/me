<script>
  import { Router, Route, Link } from 'svelte-navigator';
  import ContentForm from './components/ContentForm.svelte';
  import SearchBar from './components/SearchBar.svelte';
  import ContentList from './components/ContentList.svelte';
  
  let contents = [];
  let searchResults = [];
  let isLoading = true;
  let error = null;
  
  async function fetchContents() {
    try {
      isLoading = true;
      const response = await fetch('http://localhost:8080/api/content');
      if (response.ok) {
        contents = await response.json();
      } else {
        error = 'Failed to fetch contents';
      }
    } catch (err) {
      error = err.message;
    } finally {
      isLoading = false;
    }
  }
  
  async function handleSearch(event) {
    const query = event.detail;
    try {
      isLoading = true;
      const response = await fetch('http://localhost:8080/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          query,
          semantic: true,
          limit: 10
        })
      });
      
      if (response.ok) {
        searchResults = await response.json();
      } else {
        error = 'Search failed';
        searchResults = [];
      }
    } catch (err) {
      error = err.message;
      searchResults = [];
    } finally {
      isLoading = false;
    }
  }
  
  // Load contents when the component mounts
  fetchContents();
</script>

<Router>
  <main>
    <header>
      <h1>Personal Knowledge Base</h1>
      <div class="actions">
        <SearchBar on:search={handleSearch} />
        <Link to="/new" class="new-button">Add New</Link>
      </div>
    </header>
    
    {#if error}
      <div class="error">{error}</div>
    {/if}
    
    <div class="content">
      <Route path="/" primary={false}>
        {#if isLoading}
          <div class="loading">Loading...</div>
        {:else}
          <ContentList contents={searchResults.length ? searchResults : contents} />
        {/if}
      </Route>
      
      <Route path="/new" primary={false}>
        <ContentForm />
      </Route>
      
      <Route path="/edit/:id" primary={false} let:params>
        <ContentForm id={params.id} />
      </Route>
    </div>
  </main>
</Router>

<style>
  main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  }
  
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    flex-wrap: wrap;
    gap: 1rem;
  }
  
  h1 {
    margin: 0;
    color: #333;
  }
  
  .actions {
    display: flex;
    gap: 1rem;
    align-items: center;
    flex-wrap: wrap;
  }
  
  .content {
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
  }
  
  .new-button {
    background-color: #4caf50;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    text-decoration: none;
    font-weight: bold;
  }
  
  .error {
    background-color: #f8d7da;
    color: #721c24;
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
  
  .loading {
    text-align: center;
    padding: 2rem;
    color: #666;
  }
</style>
