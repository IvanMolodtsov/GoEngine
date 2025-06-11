package sdl

// #cgo CFLAGS: -I${SRCDIR}/include/SDL
// #cgo LDFLAGS: -L${SRCDIR}/libs -lSDL3
// #include "SDL_gpu.h"
import "C"

type GPUShaderFormat = C.SDL_GPUShaderFormat

const (
	INVALID  GPUShaderFormat = 0
	PRIVATE  GPUShaderFormat = 1 << 0 // Shaders for NDA'd platforms.
	SPIRV    GPUShaderFormat = 1 << 1 // SPIR-V shaders for Vulkan.
	DXBC     GPUShaderFormat = 1 << 2 // DXBC SM5_1 shaders for D3D12.
	DXIL     GPUShaderFormat = 1 << 3 // DXIL SM6_0 shaders for D3D12.
	MSL      GPUShaderFormat = 1 << 4 // MSL shaders for Metal
	METALLIB GPUShaderFormat = 1 << 5 // Precompiled metallib shaders for Metal.
)

type GPUDevice = C.SDL_GPUDevice

func CreateGPUDevice(format GPUShaderFormat, debug bool, name string) (*GPUDevice, error) {
	var cname *C.char
	if name == "" {
		cname = nil
	} else {
		cname = C.CString(name)
		defer Free(cname)
	}

	device := C.SDL_CreateGPUDevice(format, C.bool(debug), cname)
	if device == nil {
		return nil, GetError()
	}
	return device, nil
}

func ClaimWindowForGPUDevice(gpu *GPUDevice, window *Window) error {
	ok := C.SDL_ClaimWindowForGPUDevice(gpu, window)
	if !ok {
		return GetError()
	}
	return nil
}

type GPUCommandBuffer = C.SDL_GPUCommandBuffer

func AcquireGPUCommandBuffer(gpu *GPUDevice) (*GPUCommandBuffer, error) {
	cmd := C.SDL_AcquireGPUCommandBuffer(gpu)
	if cmd == nil {
		return nil, GetError()
	}
	return cmd, nil
}

func SubmitGPUCommandBuffer(cmdBuffer *GPUCommandBuffer) error {
	ok := C.SDL_SubmitGPUCommandBuffer(cmdBuffer)
	if !ok {
		return GetError()
	}
	return nil
}

type GPUTexture = C.SDL_GPUTexture

func WaitAndAcquireGPUSwapchainTexture(cmdBuffer *GPUCommandBuffer, window *Window, texture **GPUTexture) error {
	ok := C.SDL_WaitAndAcquireGPUSwapchainTexture(cmdBuffer, window, texture, nil, nil)
	if !ok {
		return GetError()
	}
	return nil
}

type GPULoadOp = C.SDL_GPULoadOp

const (
	GPU_LOADOP_LOAD      GPULoadOp = iota // The previous contents of the texture will be preserved.
	GPU_LOADOP_CLEAR                      // The contents of the texture will be cleared to a color.
	GPU_LOADOP_DONT_CARE                  // The previous contents of the texture need not be preserved. The contents will be undefined.
)

type GPUStoreOp = C.SDL_GPUStoreOp

const (
	GPU_STOREOP_STORE             GPUStoreOp = iota // The contents generated during the render pass will be written to memory.
	GPU_STOREOP_DONT_CARE                           // The contents generated during the render pass are not needed and may be discarded. The contents will be undefined.
	GPU_STOREOP_RESOLVE                             // The multisample contents generated during the render pass will be resolved to a non-multisample texture. The contents in the multisample texture may then be discarded and will be undefined.
	GPU_STOREOP_RESOLVE_AND_STORE                   // The multisample contents generated during the render pass will be resolved to a non-multisample texture. The contents in the multisample texture will be written to memory.
)

type GPUColorTargetInfo struct {
	Texture             *GPUTexture // The texture that will be used as a color target by a render pass.
	MipLevel            uint32      // The mip level to use as a color target.
	LayerOrDepthPlane   uint32      // The layer index or depth plane to use as a color target. This value is treated as a layer index on 2D array and cube textures, and as a depth plane on 3D textures.
	ClearColor          FColor      // The color to clear the color target to at the start of the render pass. Ignored if SDL_GPU_LOADOP_CLEAR is not used.
	LoadOp              GPULoadOp   // What is done with the contents of the color target at the beginning of the render pass.
	StoreOp             GPUStoreOp  // What is done with the results of the render pass.
	ResolveTexture      *GPUTexture // The texture that will receive the results of a multisample resolve operation. Ignored if a RESOLVE* store_op is not used.
	ResolveMipLevel     uint32      // The mip level of the resolve texture to use for the resolve operation. Ignored if a RESOLVE* store_op is not used.
	ResolveLayer        uint32      // The layer index of the resolve texture to use for the resolve operation. Ignored if a RESOLVE* store_op is not used.
	Cycle               bool        // true cycles the texture if the texture is bound and load_op is not LOAD
	CycleResolveTexture bool        // true cycles the resolve texture if the resolve texture is bound. Ignored if a RESOLVE* store_op is not used.
	padding1            C.Uint8
	padding2            C.Uint8
	_                   [4]byte
}

func (info *GPUColorTargetInfo) cPtr() *C.SDL_GPUColorTargetInfo {
	return CastPointer[C.SDL_GPUColorTargetInfo](info)
}

type GPUDepthStencilTargetInfo = C.SDL_GPUDepthStencilTargetInfo

type GPURenderPass = C.SDL_GPURenderPass

func BeginGPURenderPass(cmdBuffer *GPUCommandBuffer, colorTargetInfo *GPUColorTargetInfo, numberOfTargets uint32, stencilTargetInfo *GPUDepthStencilTargetInfo) (*GPURenderPass, error) {
	renderPass := C.SDL_BeginGPURenderPass(cmdBuffer, colorTargetInfo.cPtr(), C.Uint32(numberOfTargets), stencilTargetInfo)
	if renderPass == nil {
		return nil, GetError()
	}
	return renderPass, nil
}

func EndGPURenderPass(renderPass *GPURenderPass) {
	C.SDL_EndGPURenderPass(renderPass)
}

// An opaque handle representing a compiled shader object.
type GPUShader = C.SDL_GPUShader

// Specifies the format of a vertex attribute.
type GPUVertexElementFormat = C.SDL_GPUVertexElementFormat

const (
	SDL_GPU_VERTEXELEMENTFORMAT_INVALID GPUVertexElementFormat = iota

	/* 32-bit Signed Integers */
	SDL_GPU_VERTEXELEMENTFORMAT_INT
	SDL_GPU_VERTEXELEMENTFORMAT_INT2
	SDL_GPU_VERTEXELEMENTFORMAT_INT3
	SDL_GPU_VERTEXELEMENTFORMAT_INT4

	/* 32-bit Unsigned Integers */
	SDL_GPU_VERTEXELEMENTFORMAT_UINT
	SDL_GPU_VERTEXELEMENTFORMAT_UINT2
	SDL_GPU_VERTEXELEMENTFORMAT_UINT3
	SDL_GPU_VERTEXELEMENTFORMAT_UINT4

	/* 32-bit Floats */
	SDL_GPU_VERTEXELEMENTFORMAT_FLOAT
	SDL_GPU_VERTEXELEMENTFORMAT_FLOAT2
	SDL_GPU_VERTEXELEMENTFORMAT_FLOAT3
	SDL_GPU_VERTEXELEMENTFORMAT_FLOAT4

	/* 8-bit Signed Integers */
	SDL_GPU_VERTEXELEMENTFORMAT_BYTE2
	SDL_GPU_VERTEXELEMENTFORMAT_BYTE4

	/* 8-bit Unsigned Integers */
	SDL_GPU_VERTEXELEMENTFORMAT_UBYTE2
	SDL_GPU_VERTEXELEMENTFORMAT_UBYTE4

	/* 8-bit Signed Normalized */
	SDL_GPU_VERTEXELEMENTFORMAT_BYTE2_NORM
	SDL_GPU_VERTEXELEMENTFORMAT_BYTE4_NORM

	/* 8-bit Unsigned Normalized */
	SDL_GPU_VERTEXELEMENTFORMAT_UBYTE2_NORM
	SDL_GPU_VERTEXELEMENTFORMAT_UBYTE4_NORM

	/* 16-bit Signed Integers */
	SDL_GPU_VERTEXELEMENTFORMAT_SHORT2
	SDL_GPU_VERTEXELEMENTFORMAT_SHORT4

	/* 16-bit Unsigned Integers */
	SDL_GPU_VERTEXELEMENTFORMAT_USHORT2
	SDL_GPU_VERTEXELEMENTFORMAT_USHORT4

	/* 16-bit Signed Normalized */
	SDL_GPU_VERTEXELEMENTFORMAT_SHORT2_NORM
	SDL_GPU_VERTEXELEMENTFORMAT_SHORT4_NORM

	/* 16-bit Unsigned Normalized */
	SDL_GPU_VERTEXELEMENTFORMAT_USHORT2_NORM
	SDL_GPU_VERTEXELEMENTFORMAT_USHORT4_NORM

	/* 16-bit Floats */
	SDL_GPU_VERTEXELEMENTFORMAT_HALF2
	SDL_GPU_VERTEXELEMENTFORMAT_HALF4
)

/*
  Specifies the rate at which vertex attributes are pulled from buffers.
*/
type GPUVertexInputRate = C.SDL_GPUVertexInputRate

const (
	GPU_VERTEXINPUTRATE_VERTEX   GPUVertexInputRate = iota /* Attribute addressing is a function of the vertex index. */
	GPU_VERTEXINPUTRATE_INSTANCE                           /* Attribute addressing is a function of the instance index. */
)

/*
  A structure specifying the parameters of vertex buffers used in a graphics
  pipeline.

  When you call SDL_BindGPUVertexBuffers, you specify the binding slots of
  the vertex buffers. For example if you called SDL_BindGPUVertexBuffers with
  a first_slot of 2 and num_bindings of 3, the binding slots 2, 3, 4 would be
  used by the vertex buffers you pass in.

  Vertex attributes are linked to buffers via the buffer_slot field of
  SDL_GPUVertexAttribute. For example, if an attribute has a buffer_slot of
  0, then that attribute belongs to the vertex buffer bound at slot 0.
*/
type GPUVertexBufferDescription struct {
	Slot             uint32             //The binding slot of the vertex buffer.
	Pitch            uint32             //The byte pitch between consecutive elements of the vertex buffer.
	InputRate        GPUVertexInputRate //Whether attribute addressing is a function of the vertex index or instance index.
	InstanceStepRate uint32             //Reserved for future use. Must be set to 0.
}

/*
  A structure specifying a vertex attribute.
  All vertex attribute locations provided to an SDL_GPUVertexInputState must be unique.
*/
type GPUVertexAttribute struct {
	Location    uint32                 // The shader input location index.
	Buffer_slot uint32                 // The binding slot of the associated vertex buffer.
	Format      GPUVertexElementFormat // The size and type of the attribute data.
	Offset      uint32                 // The byte offset of this attribute relative to the start of the vertex element.
}

/*
  A structure specifying the parameters of a graphics pipeline vertex input state.
*/
type GPUVertexInputState struct {
	VertexBufferDescriptions *GPUVertexBufferDescription // A pointer to an array of vertex buffer descriptions.
	NumVertexBuffers         uint32                      // The number of vertex buffer descriptions in the above array.
	VertexAttributes         *GPUVertexAttribute         //  A pointer to an array of vertex attribute descriptions.
	NumVertexAttributes      uint32                      // The number of vertex attribute descriptions in the above array.
	_                        [4]byte
}

/*
  Specifies the primitive topology of a graphics pipeline.

  If you are using POINTLIST you must include a point size output in the
  vertex shader.

  - For HLSL compiling to SPIRV you must decorate a float output with [[vk::builtin("PointSize")]].
  - For GLSL you must set the gl_PointSize builtin.
  - For MSL you must include a float output with the [[point_size]] decorator.

  Note that sized point topology is totally unsupported on D3D12. Any size
  other than 1 will be ignored. In general, you should avoid using point
	topology for both compatibility and performance reasons. You WILL regret
  using it.

*/
type GPUPrimitiveType = C.SDL_GPUPrimitiveType

const (
	GPU_PRIMITIVETYPE_TRIANGLELIST  GPUPrimitiveType = iota // A series of separate triangles.
	GPU_PRIMITIVETYPE_TRIANGLESTRIP                         // A series of connected triangles.
	GPU_PRIMITIVETYPE_LINELIST                              // A series of separate lines.
	GPU_PRIMITIVETYPE_LINESTRIP                             // A series of connected lines.
	GPU_PRIMITIVETYPE_POINTLIST                             // A series of separate points.
)

// Specifies the fill mode of the graphics pipeline.
type GPUFillMode = C.SDL_GPUFillMode

const (
	GPU_FILLMODE_FILL GPUFillMode = iota // Polygons will be rendered via rasterization.
	GPU_FILLMODE_LINE                    // Polygon edges will be drawn as line segments.
)

// Specifies the facing direction in which triangle faces will be culled.
type GPUCullMode = C.SDL_GPUCullMode

const (
	GPU_CULLMODE_NONE  GPUCullMode = iota /* No triangles are culled. */
	GPU_CULLMODE_FRONT                    /* Front-facing triangles are culled. */
	GPU_CULLMODE_BACK                     /* Back-facing triangles are culled. */
)

/*
  Specifies the vertex winding that will cause a triangle to be determined to
  be front-facing.
*/
type GPUFrontFace = C.SDL_GPUFrontFace

const (
	GPU_FRONTFACE_COUNTER_CLOCKWISE GPUFrontFace = iota // A triangle with counter-clockwise vertex winding will be considered front-facing.
	GPU_FRONTFACE_CLOCKWISE                             // A triangle with clockwise vertex winding will be considered front-facing.
)

/*
  A structure specifying the parameters of the graphics pipeline rasterizer
  state.

  Note that SDL_GPU_FILLMODE_LINE is not supported on many Android devices.
  For those devices, the fill mode will automatically fall back to FILL.

  Also note that the D3D12 driver will enable depth clamping even if
  enable_depth_clip is true. If you need this clamp+clip behavior, consider
  enabling depth clip and then manually clamping depth in your fragment
  shaders on Metal and Vulkan.
*/
type GPURasterizerState struct {
	FillMode                GPUFillMode  // Whether polygons will be filled in or drawn as lines.
	CullMode                GPUCullMode  // The facing direction in which triangles will be culled.
	FrontFace               GPUFrontFace // The vertex winding that will cause a triangle to be determined as front-facing.
	DepthBiasConstantFactor float32      // A scalar factor controlling the depth value added to each fragment.
	DepthBiasClamp          float32      // The maximum depth bias of a fragment.
	DepthBiasSlopeFactor    float32      // A scalar factor applied to a fragment's slope in depth calculations.
	EnableDepthBias         bool         // true to bias fragment depth values.
	EnableDepthClip         bool         // true to enable depth clip, false to enable depth clamp.
	padding1                uint8
	padding2                uint8
}

/*
	Specifies the sample count of a texture.

	Used in multisampling. Note that this value only applies when the texture
	is used as a render target.
*/
type GPUSampleCount = C.SDL_GPUSampleCount

const (
	GPU_SAMPLECOUNT_1 GPUSampleCount = iota /* No multisampling. */
	GPU_SAMPLECOUNT_2                       /* MSAA 2x */
	GPU_SAMPLECOUNT_4                       /* MSAA 4x */
	GPU_SAMPLECOUNT_8                       /* MSAA 8x */
)

/*
  A structure specifying the parameters of the graphics pipeline multisample
  state.
*/
type GPUMultisampleState struct {
	SampleCount GPUSampleCount // The number of samples to be used in rasterization.
	SampleMask  uint32         //Reserved for future use. Must be set to 0.
	EnableMask  bool           //Reserved for future use. Must be set to false.
	padding1    int8
	padding2    int8
	padding3    int8
}

/*
  Specifies a comparison operator for depth, stencil and sampler operations.
*/
type GPUCompareOp = C.SDL_GPUCompareOp

const (
	GPU_COMPAREOP_INVALID          GPUCompareOp = iota
	GPU_COMPAREOP_NEVER                         /* The comparison always evaluates false. */
	GPU_COMPAREOP_LESS                          /* The comparison evaluates reference < test. */
	GPU_COMPAREOP_EQUAL                         /* The comparison evaluates reference == test. */
	GPU_COMPAREOP_LESS_OR_EQUAL                 /* The comparison evaluates reference <= test. */
	GPU_COMPAREOP_GREATER                       /* The comparison evaluates reference > test. */
	GPU_COMPAREOP_NOT_EQUAL                     /* The comparison evaluates reference != test. */
	GPU_COMPAREOP_GREATER_OR_EQUAL              /* The comparison evaluates reference >= test. */
	GPU_COMPAREOP_ALWAYS                        /* The comparison always evaluates true. */
)

/*
  Specifies what happens to a stored stencil value if stencil tests fail or
*/
type GPUStencilOp = C.SDL_GPUStencilOp

const (
	GPU_STENCILOP_INVALID             GPUStencilOp = iota
	GPU_STENCILOP_KEEP                             /* Keeps the current value. */
	GPU_STENCILOP_ZERO                             /* Sets the value to 0. */
	GPU_STENCILOP_REPLACE                          /* Sets the value to reference. */
	GPU_STENCILOP_INCREMENT_AND_CLAMP              /* Increments the current value and clamps to the maximum value. */
	GPU_STENCILOP_DECREMENT_AND_CLAMP              /* Decrements the current value and clamps to 0. */
	GPU_STENCILOP_INVERT                           /* Bitwise-inverts the current value. */
	GPU_STENCILOP_INCREMENT_AND_WRAP               /* Increments the current value and wraps back to 0. */
	GPU_STENCILOP_DECREMENT_AND_WRAP               /* Decrements the current value and wraps to the maximum value. */
)

//A structure specifying the stencil operation state of a graphics pipeline.
type GPUStencilOpState struct {
	FailOp      GPUStencilOp // The action performed on samples that fail the stencil test.
	PassOp      GPUStencilOp // The action performed on samples that pass the depth and stencil tests.
	DepthFailOp GPUStencilOp // The action performed on samples that pass the stencil test and fail the depth test.
	CompareOp   GPUCompareOp // The comparison operator used in the stencil test.
}

// A structure specifying the parameters of the graphics pipeline depth stencil state.
type GPUDepthStencilState struct {
	CompareOp         GPUCompareOp      //The comparison operator used for depth testing.
	BackStencilState  GPUStencilOpState // The stencil op state for back-facing triangles.
	FrontStencilState GPUStencilOpState //The stencil op state for front-facing triangles.
	CompareMask       uint8             //Selects the bits of the stencil values participating in the stencil test.
	WriteMask         uint8             //Selects the bits of the stencil values updated by the stencil test.
	EnableDepthTest   bool              //true enables the depth test.
	EnableDepthWrite  bool              //true enables depth writes. Depth writes are always disabled when EnableDepthTest is false.
	EnableStencilTest bool              //true enables the stencil test.
	padding1          uint8
	padding2          uint8
	padding3          uint8
}

/*
  Specifies the pixel format of a texture.

  Texture format support varies depending on driver, hardware, and usage
  flags. In general, you should use SDL_GPUTextureSupportsFormat to query if
  a format is supported before using it. However, there are a few guaranteed
  formats.

  FIXME: Check universal support for 32-bit component formats FIXME: Check
  universal support for SIMULTANEOUS_READ_WRITE

  For SAMPLER usage, the following formats are universally supported:
		- R8G8B8A8_UNORM
		- B8G8R8A8_UNORM
		- R8_UNORM
		- R8_SNORM
		- R8G8_UNORM
		- R8G8_SNORM
		- R8G8B8A8_SNORM
		- R16_FLOAT
		- R16G16_FLOAT
		- R16G16B16A16_FLOAT
		- R32_FLOAT
		- R32G32_FLOAT
		- R32G32B32A32_FLOAT
		- R11G11B10_UFLOAT
		- R8G8B8A8_UNORM_SRGB
		- B8G8R8A8_UNORM_SRGB
		- D16_UNORM


  For COLOR_TARGET usage, the following formats are universally supported:
		- R8G8B8A8_UNORM
		- B8G8R8A8_UNORM
		- R8_UNORM
		- R16_FLOAT
		- R16G16_FLOAT
		- R16G16B16A16_FLOAT
		- R32_FLOAT
		- R32G32_FLOAT
		- R32G32B32A32_FLOAT
		- R8_UINT
		- R8G8_UINT
		- R8G8B8A8_UINT
		- R16_UINT
		- R16G16_UINT
		- R16G16B16A16_UINT
		- R8_INT
		- R8G8_INT
		- R8G8B8A8_INT
		- R16_INT
		- R16G16_INT
		- R16G16B16A16_INT
		- R8G8B8A8_UNORM_SRGB
		- B8G8R8A8_UNORM_SRGB


  For STORAGE usages, the following formats are universally supported:
		- R8G8B8A8_UNORM
		- R8G8B8A8_SNORM
		- R16G16B16A16_FLOAT
		- R32_FLOAT
		- R32G32_FLOAT
		- R32G32B32A32_FLOAT
		- R8G8B8A8_UINT
		- R16G16B16A16_UINT
		- R8G8B8A8_INT
		- R16G16B16A16_INT

  For DEPTH_STENCIL_TARGET usage, the following formats are universally
  supported:
		- D16_UNORM
		- Either (but not necessarily both!) D24_UNORM or D32_FLOAT
		- Either (but not necessarily both!) D24_UNORM_S8_UINT or D32_FLOAT_S8_UINT

  Unless D16_UNORM is sufficient for your purposes, always check which of
  D24/D32 is supported before creating a depth-stencil texture!

*/
type GPUTextureFormat = C.SDL_GPUTextureFormat

const (
	SDL_GPU_TEXTUREFORMAT_INVALID GPUTextureFormat = iota

	/* Unsigned Normalized Float Color Formats */
	SDL_GPU_TEXTUREFORMAT_A8_UNORM
	SDL_GPU_TEXTUREFORMAT_R8_UNORM
	SDL_GPU_TEXTUREFORMAT_R8G8_UNORM
	SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UNORM
	SDL_GPU_TEXTUREFORMAT_R16_UNORM
	SDL_GPU_TEXTUREFORMAT_R16G16_UNORM
	SDL_GPU_TEXTUREFORMAT_R16G16B16A16_UNORM
	SDL_GPU_TEXTUREFORMAT_R10G10B10A2_UNORM
	SDL_GPU_TEXTUREFORMAT_B5G6R5_UNORM
	SDL_GPU_TEXTUREFORMAT_B5G5R5A1_UNORM
	SDL_GPU_TEXTUREFORMAT_B4G4R4A4_UNORM
	SDL_GPU_TEXTUREFORMAT_B8G8R8A8_UNORM
	/* Compressed Unsigned Normalized Float Color Formats */
	SDL_GPU_TEXTUREFORMAT_BC1_RGBA_UNORM
	SDL_GPU_TEXTUREFORMAT_BC2_RGBA_UNORM
	SDL_GPU_TEXTUREFORMAT_BC3_RGBA_UNORM
	SDL_GPU_TEXTUREFORMAT_BC4_R_UNORM
	SDL_GPU_TEXTUREFORMAT_BC5_RG_UNORM
	SDL_GPU_TEXTUREFORMAT_BC7_RGBA_UNORM
	/* Compressed Signed Float Color Formats */
	SDL_GPU_TEXTUREFORMAT_BC6H_RGB_FLOAT
	/* Compressed Unsigned Float Color Formats */
	SDL_GPU_TEXTUREFORMAT_BC6H_RGB_UFLOAT
	/* Signed Normalized Float Color Formats  */
	SDL_GPU_TEXTUREFORMAT_R8_SNORM
	SDL_GPU_TEXTUREFORMAT_R8G8_SNORM
	SDL_GPU_TEXTUREFORMAT_R8G8B8A8_SNORM
	SDL_GPU_TEXTUREFORMAT_R16_SNORM
	SDL_GPU_TEXTUREFORMAT_R16G16_SNORM
	SDL_GPU_TEXTUREFORMAT_R16G16B16A16_SNORM
	/* Signed Float Color Formats */
	SDL_GPU_TEXTUREFORMAT_R16_FLOAT
	SDL_GPU_TEXTUREFORMAT_R16G16_FLOAT
	SDL_GPU_TEXTUREFORMAT_R16G16B16A16_FLOAT
	SDL_GPU_TEXTUREFORMAT_R32_FLOAT
	SDL_GPU_TEXTUREFORMAT_R32G32_FLOAT
	SDL_GPU_TEXTUREFORMAT_R32G32B32A32_FLOAT
	/* Unsigned Float Color Formats */
	SDL_GPU_TEXTUREFORMAT_R11G11B10_UFLOAT
	/* Unsigned Integer Color Formats */
	SDL_GPU_TEXTUREFORMAT_R8_UINT
	SDL_GPU_TEXTUREFORMAT_R8G8_UINT
	SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UINT
	SDL_GPU_TEXTUREFORMAT_R16_UINT
	SDL_GPU_TEXTUREFORMAT_R16G16_UINT
	SDL_GPU_TEXTUREFORMAT_R16G16B16A16_UINT
	SDL_GPU_TEXTUREFORMAT_R32_UINT
	SDL_GPU_TEXTUREFORMAT_R32G32_UINT
	SDL_GPU_TEXTUREFORMAT_R32G32B32A32_UINT
	/* Signed Integer Color Formats */
	SDL_GPU_TEXTUREFORMAT_R8_INT
	SDL_GPU_TEXTUREFORMAT_R8G8_INT
	SDL_GPU_TEXTUREFORMAT_R8G8B8A8_INT
	SDL_GPU_TEXTUREFORMAT_R16_INT
	SDL_GPU_TEXTUREFORMAT_R16G16_INT
	SDL_GPU_TEXTUREFORMAT_R16G16B16A16_INT
	SDL_GPU_TEXTUREFORMAT_R32_INT
	SDL_GPU_TEXTUREFORMAT_R32G32_INT
	SDL_GPU_TEXTUREFORMAT_R32G32B32A32_INT
	/* SRGB Unsigned Normalized Color Formats */
	SDL_GPU_TEXTUREFORMAT_R8G8B8A8_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_B8G8R8A8_UNORM_SRGB
	/* Compressed SRGB Unsigned Normalized Color Formats */
	SDL_GPU_TEXTUREFORMAT_BC1_RGBA_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_BC2_RGBA_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_BC3_RGBA_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_BC7_RGBA_UNORM_SRGB
	/* Depth Formats */
	SDL_GPU_TEXTUREFORMAT_D16_UNORM
	SDL_GPU_TEXTUREFORMAT_D24_UNORM
	SDL_GPU_TEXTUREFORMAT_D32_FLOAT
	SDL_GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT
	SDL_GPU_TEXTUREFORMAT_D32_FLOAT_S8_UINT
	/* Compressed ASTC Normalized Float Color Formats*/
	SDL_GPU_TEXTUREFORMAT_ASTC_4x4_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_5x4_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_5x5_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_6x5_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_6x6_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_8x5_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_8x6_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_8x8_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_10x5_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_10x6_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_10x8_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_10x10_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_12x10_UNORM
	SDL_GPU_TEXTUREFORMAT_ASTC_12x12_UNORM
	/* Compressed SRGB ASTC Normalized Float Color Formats*/
	SDL_GPU_TEXTUREFORMAT_ASTC_4x4_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_5x4_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_5x5_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_6x5_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_6x6_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_8x5_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_8x6_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_8x8_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_10x5_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_10x6_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_10x8_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_10x10_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_12x10_UNORM_SRGB
	SDL_GPU_TEXTUREFORMAT_ASTC_12x12_UNORM_SRGB
	/* Compressed ASTC Signed Float Color Formats*/
	SDL_GPU_TEXTUREFORMAT_ASTC_4x4_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_5x4_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_5x5_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_6x5_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_6x6_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_8x5_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_8x6_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_8x8_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_10x5_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_10x6_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_10x8_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_10x10_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_12x10_FLOAT
	SDL_GPU_TEXTUREFORMAT_ASTC_12x12_FLOAT
)

type GPUBlendFactor = C.SDL_GPUBlendFactor
type GPUBlendOp = C.SDL_GPUBlendOp
type GPUColorComponentFlags = C.SDL_GPUColorComponentFlags

type GPUColorTargetBlendState struct {
	src_color_blendfactor   GPUBlendFactor
	dst_color_blendfactor   GPUBlendFactor
	color_blend_op          GPUBlendOp
	src_alpha_blendfactor   GPUBlendFactor
	dst_alpha_blendfactor   GPUBlendFactor
	alpha_blend_op          GPUBlendOp
	color_write_mask        GPUColorComponentFlags
	enable_blend            bool
	enable_color_write_mask bool
	padding1                uint8
	padding2                uint8
	_                       [3]byte
}

type GPUColorTargetDescription struct {
	format      GPUTextureFormat
	blend_state GPUColorTargetBlendState
}

type GPUGraphicsPipelineTargetInfo struct {
	color_target_descriptions *GPUColorTargetDescription
	num_color_targets         uint32
	depth_stencil_format      GPUTextureFormat
	has_depth_stencil_target  bool
	padding1                  uint8
	padding2                  uint8
	padding3                  uint8
	_                         [4]byte
}

type GPUGraphicsPipelineCreateInfo struct {
	VertexShader      *GPUShader
	FragmentShader    *GPUShader
	VertexInputState  GPUVertexInputState
	PrimitiveType     GPUPrimitiveType
	RasterizerState   GPURasterizerState
	MultisampleState  GPUMultisampleState
	DepthStencilState GPUDepthStencilState
	TargetInfo        GPUGraphicsPipelineTargetInfo //Formats and blend modes for the render targets of the graphics pipeline.
	Props             PropertiesID                  // A properties ID for extensions. Should be 0 if no extensions are needed.
	_                 [4]byte
}

type GPUGraphicsPipeline = C.SDL_GPUGraphicsPipeline

func CreateGPUGraphicsPipeline(gpu *GPUDevice, createInfo *GPUGraphicsPipelineCreateInfo) (*GPUGraphicsPipeline, error) {
	pipe := C.SDL_CreateGPUGraphicsPipeline(gpu, CastPointer[C.SDL_GPUGraphicsPipelineCreateInfo](createInfo))
	if pipe == nil {
		return nil, GetError()
	}
	return pipe, nil
}
