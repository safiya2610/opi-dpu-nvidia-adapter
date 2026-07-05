# OPI DPU NVIDIA Adapter

This project is a small proof-of-concept operator that adds NVIDIA support to an OPI-based DPU workflow.
It uses an adapter layer to translate an OPI DPUCluster into NVIDIA DPF resources such as DPUSet and DPUService.

## What This Project Does

- Accepts an OPI DPUCluster resource as input.
- Translates it into NVIDIA DPF child resources.
- Reconciles those resources with a controller-runtime reconciler.
- Updates the parent OPI status based on child resource readiness.

## Why This Design

The main idea is to keep the OPI layer vendor-neutral while reusing NVIDIA's existing DPF workflow instead of rebuilding everything from scratch.
This makes the code cleaner and easier to extend for future vendors.

## Project Structure

- api/v1alpha1: OPI API types
- controllers: reconciliation logic
- pkg/adapter: translation logic for NVIDIA
- examples: demo program showing the translation flow

## How It Works

1. A user creates an OPI DPUCluster resource.
2. The reconciler detects the vendor as NVIDIA.
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

The output shows the translation working successfully.
