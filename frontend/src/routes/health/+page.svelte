<script lang="ts">
    // 'data' enthält alles, was die +page.server.ts zurückgegeben hat
    let { data } = $props(); 
</script>

<div class="p-8 font-sans">
    <h1 class="text-2xl font-bold mb-4">System Status</h1>
    
    <div class="bg-gray-100 p-4 rounded-lg shadow-md inline-block">
        <p><strong>Status:</strong> 
            <span class="text-green-600">{data.backendData.status}</span>
        </p>
        <p><strong>Nachricht:</strong> {data.backendData.message}</p>
    </div>
</div>