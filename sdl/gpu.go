package sdl

// #cgo CFLAGS: -I${SRCDIR}/include/SDL
// #cgo LDFLAGS: -L${SRCDIR}/libs -lSDL3
// #include "SDL.h"
import "C"

type GPUShaderFormat = C.SDL_GPUShaderFormat

const (
	INVALID  GPUShaderFormat = 0
	PRIVATE  GPUShaderFormat = 1 << 0
	SPIRV    GPUShaderFormat = 1 << 1
	DXBC     GPUShaderFormat = 1 << 2
	DXIL     GPUShaderFormat = 1 << 3
	MSL      GPUShaderFormat = 1 << 4
	METALLIB GPUShaderFormat = 1 << 5
)

type GPUDevice = C.SDL_GPUDevice

func CreateGPUDevice(format GPUShaderFormat, debug bool, name string) (*GPUDevice, error) {
	var cname *C.char
	if name == "" {
		cname = nil
	} else {
		cname = C.CString(name)
	}
	defer Free(cname)
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
	GPU_LOADOP_LOAD      GPULoadOp = iota /**< The previous contents of the texture will be preserved. */
	GPU_LOADOP_CLEAR     GPULoadOp = iota /**< The contents of the texture will be cleared to a color. */
	GPU_LOADOP_DONT_CARE GPULoadOp = iota /**< The previous contents of the texture need not be preserved. The contents will be undefined. */
)

type GPUStoreOp = C.SDL_GPUStoreOp

const (
	GPU_STOREOP_STORE             GPUStoreOp = iota /**< The contents generated during the render pass will be written to memory. */
	GPU_STOREOP_DONT_CARE         GPUStoreOp = iota /**< The contents generated during the render pass are not needed and may be discarded. The contents will be undefined. */
	GPU_STOREOP_RESOLVE           GPUStoreOp = iota /**< The multisample contents generated during the render pass will be resolved to a non-multisample texture. The contents in the multisample texture may then be discarded and will be undefined. */
	GPU_STOREOP_RESOLVE_AND_STORE GPUStoreOp = iota /**< The multisample contents generated during the render pass will be resolved to a non-multisample texture. The contents in the multisample texture will be written to memory. */
)

type GPUColorTargetInfo struct {
	Texture             *GPUTexture /**< The texture that will be used as a color target by a render pass. */
	MipLevel            uint32      /**< The mip level to use as a color target. */
	LayerOrDepthPlane   uint32      /**< The layer index or depth plane to use as a color target. This value is treated as a layer index on 2D array and cube textures, and as a depth plane on 3D textures. */
	ClearColor          FColor      /**< The color to clear the color target to at the start of the render pass. Ignored if SDL_GPU_LOADOP_CLEAR is not used. */
	LoadOp              GPULoadOp   /**< What is done with the contents of the color target at the beginning of the render pass. */
	StoreOp             GPUStoreOp  /**< What is done with the results of the render pass. */
	ResolveTexture      *GPUTexture /**< The texture that will receive the results of a multisample resolve operation. Ignored if a RESOLVE* store_op is not used. */
	ResolveMipLevel     uint32      /**< The mip level of the resolve texture to use for the resolve operation. Ignored if a RESOLVE* store_op is not used. */
	ResolveLayer        uint32      /**< The layer index of the resolve texture to use for the resolve operation. Ignored if a RESOLVE* store_op is not used. */
	Cycle               bool        /**< true cycles the texture if the texture is bound and load_op is not LOAD */
	CycleResolveTexture bool        /**< true cycles the resolve texture if the resolve texture is bound. Ignored if a RESOLVE* store_op is not used. */
	padding1            C.Uint8
	padding2            C.Uint8
	// _                   [4]byte
}

func (info GPUColorTargetInfo) CStruct() C.SDL_GPUColorTargetInfo {
	return C.SDL_GPUColorTargetInfo{
		texture:               info.Texture,
		mip_level:             C.Uint32(info.MipLevel),
		layer_or_depth_plane:  C.Uint32(info.LayerOrDepthPlane),
		clear_color:           info.ClearColor.CStruct(),
		load_op:               info.LoadOp,
		store_op:              info.StoreOp,
		resolve_texture:       info.ResolveTexture,
		resolve_mip_level:     C.Uint32(info.ResolveMipLevel),
		resolve_layer:         C.Uint32(info.ResolveLayer),
		cycle:                 C.bool(info.Cycle),
		cycle_resolve_texture: C.bool(info.CycleResolveTexture),
	}
}

// type GPUColorTargetInfo = C.SDL_GPUColorTargetInfo
type GPUDepthStencilTargetInfo = C.SDL_GPUDepthStencilTargetInfo

type GPURenderPass = C.SDL_GPURenderPass

func BeginGPURenderPass(cmdBuffer *GPUCommandBuffer, colorTargetInfo GPUColorTargetInfo, numberOfTargets uint32, stencilTargetInfo *GPUDepthStencilTargetInfo) (*GPURenderPass, error) {
	temp := colorTargetInfo.CStruct()

	renderPass := C.SDL_BeginGPURenderPass(cmdBuffer, &temp, C.Uint32(numberOfTargets), stencilTargetInfo)
	if renderPass == nil {
		return nil, GetError()
	}
	return renderPass, nil
}

func EndGPURenderPass(renderPass *GPURenderPass) {
	C.SDL_EndGPURenderPass(renderPass)
}
