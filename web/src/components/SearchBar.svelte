<script>
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  let query = '';
  let semantic = true;
  
  function handleSubmit() {
    if (query.trim()) {
      dispatch('search', query);
    }
  }
</script>

<div class="search-container">
  <form on:submit|preventDefault={handleSubmit}>
    <input 
      type="text" 
      placeholder="Search your knowledge base..." 
      bind:value={query}
      aria-label="Search query"
    >
    <div class="search-options">
      <label>
        <input type="checkbox" bind:checked={semantic}>
        Semantic Search
      </label>
    </div>
    <button type="submit" aria-label="Search">Search</button>
  </form>
</div>

<style>
  .search-container {
    width: 100%;
    max-width: 600px;
  }
  
  form {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    align-items: center;
  }
  
  input[type="text"] {
    flex: 1;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    min-width: 200px;
    font-size: 1rem;
  }
  
  .search-options {
    display: flex;
    align-items: center;
    font-size: 0.875rem;
    color: #555;
  }
  
  .search-options label {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    cursor: pointer;
  }
  
  button {
    background-color: #2196f3;
    color: white;
    border: none;
    padding: 0.75rem 1.25rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    font-size: 1rem;
  }
  
  button:hover {
    background-color: #1976d2;
  }
  
  @media (max-width: 600px) {
    form {
      flex-direction: column;
      align-items: stretch;
    }
    
    .search-options {
      margin: 0.5rem 0;
    }
  }
</style>
