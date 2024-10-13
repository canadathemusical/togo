brainstorming (feel free to ignore)

what if every row is a commit so, the hash and message serve as id and name
each commit has a tiny plaintext json or some machine readable file with status and notes
each branch is a tag or mode like work, personal, school, etc
and commits just get ammended?

nope that's too crazy

godo this serves as a way to list all incomplete and is an alias for godo --list
godo --add by itself opens your editor to create a new item that way  godo --add "title" creates a new item with no notes 
godo --edit <hash> will open your editor to edit an existing item godo --edit <hash> "title" will overwrite the title of an existing item
godo --done <hash> marks it as done godo --done by itself should open an interactive cli to select an item to mark as done
godo list mode <mode>  should show all items in that mode and the mode flag should be usable in --add as well
The user should be able to set a default mode which is saved to .config/godo/config.yaml or some shit
godo --delete <hash> should be the only way to delete an item and should default to soft delete, sending the item to the deleted mode, godo --delete --hard <hash> should be possible but prompt the user to be damn sure
godo --cloudsave maybe to add commit push all items to an remote git repo with some additional commands for pulling, initialization, pushing, etc

I think that's the entire ui. if you decide to read this, let me know if I missed anything useful for v1, which will probably just use sqlite and cloudsave 
