// Custom hooks for Reduxß
import { useDispatch, useSelector, type TypedUseSelectorHook } from "react-redux";
import type { AppDispatch, RootState } from "../features/store";

// Type-safe useDispatch hook
export const useAppDispatch = () => useDispatch<AppDispatch>();

// Type-safe useSelector hook
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;