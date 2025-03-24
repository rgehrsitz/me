<script>
  import { navigate } from 'svelte-navigator';
  
  export let id = null;
  
  let content = {
    type: 'note',
    title: '',
    body: '',
    source_url: '',
    file_path: '',
    tags: []
  };
  
  let loading = false;
  let error = null;
  let tagInput = '';
  
  async function fetchContent() {
    if (!id) return;
    
    try {
      loading = true;
      const response = await fetch(`http://localhost:8080/api/content/${id}`);
      if (response.ok) {
        content = await response.json();
        if (typeof content.tags === 'string') {
          content.tags = content.tags.split(',').map(tag => tag.trim());
        }
      } else {
        error = 'Failed to load content';
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
  
  async function handleSubmit() {
    try {
      loading = true;
      const url = id 
        ? `http://localhost:8080/api/content/${id}`
        : 'http://localhost:8080/api/content';
      
      const method = id ? 'PUT' : 'POST';
      
      // Process tags
      if (tagInput) {
        const newTags = tagInput.split(',').map(tag => tag.trim()).filter(tag => tag);
        content.tags = [...new Set([...content.tags, ...newTags])];
      }
      
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(content)
      });
      
      if (response.ok) {
        navigate('/');
      } else {
        const data = await response.json();
        error = data.error || 'Failed to save content';
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }
  
  function removeTag(index) {
    content.tags = content.tags.filter((_, i) => i !== index);
  }
  
  // Load content if editing
  if (id) {
    fetchContent();
  }
</script>

<div class="form-container">
  <h2>{id ? 'Edit' : 'Add'} Content</h2>
  
  {#if error}
    <div class="error">{error}</div>
  {/if}
  
  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-group">
      <label for="type">Type</label>
      <select id="type" bind:value={content.type}>
        <option value="note">Note</option>
        <option value="snippet">Code Snippet</option>
        <option value="bookmark">Bookmark</option>
        <option value="document">Document</option>
      </select>
    </div>
    
    <div class="form-group">
      <label for="title">Title</label>
      <input type="text" id="title" bind:value={content.title} required>
    </div>
    
    <div class="form-group">
      <label for="body">Content</label>
      <textarea id="body" bind:value={content.body} rows="10" required></textarea>
    </div>
    
    {#if content.type === 'bookmark'}
      <div class="form-group">
        <label for="source_url">URL</label>
        <input type="url" id="source_url" bind:value={content.source_url}>
      </div>
    {/if}
    
    {#if content.type === 'document'}
      <div class="form-group">
        <label for="file_path">File Path</label>
        <input type="text" id="file_path" bind:value={content.file_path}>
      </div>
    {/if}
    
    <div class="form-group">
      <label for="tags">Tags</label>
      <div class="tags-container">
        {#each content.tags as tag, index}
          <div class="tag">
            {tag}
            <button type="button" class="remove-tag" on:click={() => removeTag(index)}>×</button>
          </div>
        {/each}
      </div>
      <input type="text" id="tags" bind:value={tagInput} placeholder="Add tags (comma separated)">
    </div>
    
    <div class="form-actions">
      <button type="button" class="cancel-button" on:click={() => navigate('/')}>Cancel</button>
      <button type="submit" class="submit-button" disabled={loading}>
        {loading ? 'Saving...' : 'Save'}
      </button>
    </div>
  </form>
</div>

<style>
  .form-container {
    max-width: 800px;
    margin: 0 auto;
  }
  
  h2 {
    margin-top: 0;
    color: #333;
  }
  
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
    color: #555;
  }
  
  input, textarea, select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    font-family: inherit;
  }
  
  textarea {
    resize: vertical;
  }
  
  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
  }
  
  button {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    font-size: 1rem;
  }
  
  .submit-button {
    background-color: #4caf50;
    color: white;
  }
  
  .submit-button:disabled {
    background-color: #a5d6a7;
    cursor: not-allowed;
  }
  
  .cancel-button {
    background-color: #f1f1f1;
    color: #333;
  }
  
  .error {
    color: #d32f2f;
    background-color: #ffebee;
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
  }
  
  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
  
  .tag {
    background-color: #e0f2f1;
    color: #00796b;
    padding: 0.25rem 0.5rem;
    border-radius: 16px;
    display: flex;
    align-items: center;
    font-size: 0.875rem;
  }
  
  .remove-tag {
    background: none;
    border: none;
    color: #00796b;
    cursor: pointer;
    font-size: 1.25rem;
    padding: 0 0 0 0.25rem;
    line-height: 1;
  }
</style>
