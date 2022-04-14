##Electron GUI -- Application Framework:
- Electron HTML/Node.JS Interface Client.
  - AJAX/Web Sockets - Callback Event Prompts Go Server API Node
  - API Node Is Go Component/Entity System: Key/Value/Object Map
  - {Con: WebGL buffer integration with DieselFluid Framework}
  - {Con: Multiple Window Contexts for a single application, non integrated}
  - {Con: Only using Electron for Windowing/Event/JS}
  - {Pro: Configurable, CSS, Floating Task Bar, HTML 5 Animation}
  - {Pro: Complex GUI, Filesystems, Object Graphs, Embed Image/Video}
  - {Pro: GUI Completely Extensible + HTML5 + Web Socket Requests Node JS}
  - {Pro: Commit Working Notes}
  - {Pro: Don't have to implement GL based picking UI framework with event system}

  ##Current TODO:
  --Create dslfluid.com/dsl/scene package - this should be used as a non-generic structured graph representation of scene construction for the fluid simulation. I.E:
    Scene Contains (SPHMethod, Particle_Description, 3DMesheEntity Collision(Tree), 3DMeshEntity(Tree), Enumerations, Scaling, Origins, RenderAPIContext, Graphics3D, Parameterizations etc... etc...)
    -Scene (Load, Init, Render, Run)
  -- dslfluid.com/dsl/test - Test API packages composed together in functioning application for degug
  -- dslfluid.com/dsl/app/network - Open Application over Network Socket -- Define Application End Point Registry -- Define JSON Object Factory Handling


##Style Notes
 - A parent package in a module hierarchy should never be referenced by the internal module
 - High level modules should always be polymorphic
 - Module data type definitions will now be consolidated into the Types.go file for that module
 except for cases where interelated components need to exist in another file for clarity.




 ##Feature Update Status
 - Moving GLR scene rendering call routines to the Upper Module Rendersystem calls,
 - all rendering calls to the lower GLR (renderer interface) should be object instance
 - calls.
