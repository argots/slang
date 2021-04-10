# Roadmap

1. Portable Slang parser with simple syntax that is the basis for everything else
2. Portable Slang interpreter with FFI using HTTP
4. A generic package manager built on top of slang
5. Rich text/data format (rich text, inline tables, lists, etc but all static)
6. Rich text/data viewer (viewer for above with UI hierachy/views support)
7. Network Protocol Ladder generator (DSL which emits rich text/data format defined above)
8. Network Graph generator (DSL which emits rich text/data format defined above)
9. Filter/Group-by for Rich text/data format (VISIO/Figma style alternate views of data) but without a delta protocol (i.e. custom UI)
10. Static Configuration service using slang + Rich text/data viewer. (fig built on top of this)
11. Language definition of Streaming data (deltas, methods to access this in slang) 
12. Updated interpreter to use this (integreate dotchain/dot)
14. Update package manager to use this to support package changes this way.
15. Update package manager to support error reports being propagated up to authors and version upgrades.
16. Update rich text viewer => rich text editor using the streaming data protocol.
17. Integrate protocol ladders and graph generator etc into editor to support natural edit mode.
18. Update filter/group-by etc to use streaming edit over custom UI
19. Integrate editor into configuration service (editing configuration).


# TBD 

1. when does `annotations` happen? (allowing data to have annotations)
2. when to tackle the `references` problem (each data has a chain of lexical and runtime history)
3. when to tackle type systems
4. when to tackle shared service fabric/runtime?

