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
  --Polar Coordinate Mapping issues
  --Move Rendersystem to shader module component instead of the ad hoc shader compilation
  --Atmosphere HDR Texture Handling + Shader HDR Texture Handling
  --DSL App TCP Message Handling Parameter Controls
  --MOVE FLUID PARTICLES + FLUID GRIDS TO GPU GLSL SHADER LIBRARIES


##Style Notes
 - A parent package in a module hierarchy should never be referenced by the internal module
 - High level modules should always be polymorphic
 - Module data type definitions will now be consolidated into the Types.go file for that module
 except for cases where interelated components need to exist in another file for clarity.




 ##Feature Update Status
 - Moving GLR scene rendering call routines to the Upper Module Rendersystem calls,
 - all rendering calls to the lower GLR (renderer interface) should be object instance
 - calls.
