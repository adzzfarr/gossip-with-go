import { FormControl, InputLabel, MenuItem, Select, type SelectChangeEvent } from "@mui/material";

interface SortOption {
    value: string;
    label: string;
}

interface SortDropdownProps {
    value: string;
    label?: string;
    options: SortOption[];
    onChange: (value: string) => void;
}

export default function SortDropdown({
    value,
    label,
    options,
    onChange
} : SortDropdownProps) {
    const handleChange = (e: SelectChangeEvent) => {
        onChange(e.target.value);
    }

    return (
        <FormControl size="small" sx={{ minWidth: 150 }}>
            <InputLabel id="sort-dropdown-label">{label}</InputLabel>
            <Select
                labelId="sort-dropdown-label"
                id="sort-dropdown"
                value={value}
                label={label}
                onChange={handleChange}
            >
                {options.map(
                    (option) => (
                        <MenuItem key={option.value} value={option.value}>
                            {option.label}
                        </MenuItem>
                    )
                )}
            </Select>
        </FormControl>
    );
}