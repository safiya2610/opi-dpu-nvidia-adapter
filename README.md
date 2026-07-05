# OPI DPU NVIDIA Adapter

This project is a small proof-of-concept operator that adds NVIDIA support to an OPI-based DPU workflow.
It uses an adapter layer to translate an OPI DPUCluster into NVIDIA DPF resources such as DPUSet and DPUService.

## What I Did

I added NVIDIA support to the OPI operator by building a simple translation layer.
Instead of writing NVIDIA logic directly inside the OPI code, I used an adapter so the main OPI flow stays clean and vendor-neutral.
The controller now takes an OPI DPUCluster, translates it into NVIDIA resources, creates or updates them, and reflects the status back to the parent object.

## Why This Design

The main idea is to keep the OPI layer simple while reusing NVIDIA's existing DPF workflow instead of rebuilding everything from scratch.
This makes the code easier to understand and easier to extend later for other vendors.

## Project Structure

| File/Folder | Purpose |
|---|---|
| api/v1alpha1 | Defines the OPI API types used by the operator. |
| controllers | Contains the reconciliation logic for creating and updating NVIDIA resources. |
| pkg/adapter | Holds the adapter and translation code for mapping OPI input to NVIDIA DPF objects. |
| examples | Includes a demo program that shows the translation flow in a simple way. |
| main.go | Starts the controller manager and registers the reconciler. |

## How It Works

1. A user creates an OPI DPUCluster resource.
2. The reconciler checks the vendor and identifies it as NVIDIA.
3. The adapter converts the OPI object into NVIDIA DPF resources.
4. The controller creates or updates those resources.
5. Status from the NVIDIA resources is reflected back to the OPI resource.

## Demo

A simple demo is available in the examples folder.
It shows how an OPI DPUCluster is converted into NVIDIA resources.

Run the demo:

```bash
go run ./examples
```

This demo prints the input OPI spec and the generated NVIDIA DPF objects.

## Tests

Run the full test suite:

```bash
go test ./...
```

## Local Verification

The demo example was verified locally with:

```bash
go run ./examples
```

The output showed that the translation flow worked successfully.
