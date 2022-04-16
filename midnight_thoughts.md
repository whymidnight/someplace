# Midnight Thoughts

Since a Candy Machine (Batch) is space finite, _Batches_ are intended to represent as shards of inventory in the entire marketplace catalog.


We can illustrate the design like -



<Candy>0..4 -> <Batch>1 -> <Marketplace>0
<Candy>0..4 -> <Batch>2 -> <Marketplace>0
<Candy>0..4 -> <Batch>3 -> <Marketplace>0
<Candy>0..4 -> <Batch>4 -> <Marketplace>0


All <Candy> lacks the correlation of its listing meta to itself. To solve this,

we need to scrape every batch of its `Vec<<Candy>>`(err...`Vec<<ConfigLine>>`)
