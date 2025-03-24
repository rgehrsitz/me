<script>
  import { Link } from 'svelte-navigator';
  
  export let contents = [];
  
  // Format date for display
  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString();
  }
  
  // Truncate text for preview
  function truncate(text, length = 150) {
    if (!text) return '';
    return text.length > length ? text.substring(0, length) + '...' : text;
  }
  
  // Get icon for content type
  function getTypeIcon(type) {
    switch (type) {
      case 'note': return '📝';
      case 'snippet': return '💻';
      case 'bookmark': return '🔖';
      case 'document': return '📄';
      default: return '📄';
    }
  }
</script>

<div class="content-list">
  {#if contents.length === 0}
    <div class="empty-state">
      <p>No content found. Start by adding some notes, snippets, or bookmarks.</p>
      <Link to="/new" class="add-button">Add Content</Link>
    </div>
  {:else}
    <div class="list">
      {#each contents as content (content.id)}
        <div class="content-item">
          <div class="content-header">
            <span class="content-type">{getTypeIcon(content.type)} {content.type}</span>
            <span class="content-date">{formatDate(content.created_at)}</span>
          </div>
          <h3 class="content-title">{content.title}</h3>
          
          {#if content.snippet}
            <p class="content-snippet">{content.snippet}</p>
          {:else}
            <p class="content-preview">{truncate(content.body)}</p>
          {/if}
          
          {#if content.tags && content.tags.length > 0}
            <div class="content-tags">
              {#each content.tags as tag}
                <span class="tag">{tag}</span>
              {/each}
            </div>
          {/if}
          
          <div class="content-actions">
            <Link to={`/edit/${content.id}`} class="edit-button">Edit</Link>
            {#if content.type === 'bookmark' && content.source_url}
              <a href={content.source_url} target="_blank" rel="noopener noreferrer" class="view-button">Visit URL</a>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .content-list {
    width: 100%;
  }
  
  .empty-state {
    text-align: center;
    padding: 3rem 1rem;
    color: #666;
  }
  
  .add-button {
    display: inline-block;
    background-color: #4caf50;
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    text-decoration: none;
    font-weight: bold;
    margin-top: 1rem;
  }
  
  .list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }
  
  .content-item {
    border: 1px solid #eee;
    border-radius: 8px;
    padding: 1.5rem;
    background-color: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    transition: transform 0.2s, box-shadow 0.2s;
  }
  
  .content-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
  
  .content-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
    font-size: 0.875rem;
  }
  
  .content-type {
    color: #555;
    font-weight: 500;
  }
  
  .content-date {
    color: #888;
  }
  
  .content-title {
    margin: 0.5rem 0 1rem;
    color: #333;
    font-size: 1.25rem;
  }
  
  .content-preview, .content-snippet {
    color: #555;
    margin-bottom: 1rem;
    line-height: 1.5;
    font-size: 0.9375rem;
  }
  
  .content-snippet {
    background-color: #f9f9f9;
    padding: 0.75rem;
    border-radius: 4px;
    border-left: 3px solid #2196f3;
  }
  
  .content-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  
  .tag {
    background-color: #e0f2f1;
    color: #00796b;
    padding: 0.25rem 0.5rem;
    border-radius: 16px;
    font-size: 0.75rem;
  }
  
  .content-actions {
    display: flex;
    gap: 0.75rem;
    margin-top: 1rem;
  }
  
  .edit-button, .view-button {
    padding: 0.5rem 0.75rem;
    border-radius: 4px;
    text-decoration: none;
    font-size: 0.875rem;
    font-weight: 500;
  }
  
  .edit-button {
    background-color: #f1f1f1;
    color: #333;
  }
  
  .view-button {
    background-color: #e3f2fd;
    color: #1976d2;
  }
  
  @media (max-width: 600px) {
    .list {
      grid-template-columns: 1fr;
    }
  }
</style>
