import { ChevronLeft, ChevronRight, FirstPage, LastPage } from "@mui/icons-material";
import { Box, Button, IconButton, Typography } from "@mui/material";

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    disabled?: boolean;
    onPageChange: (page: number) => void;
}

export default function Pagination({
    currentPage,
    totalPages,
    disabled = false,
    onPageChange
}: PaginationProps) {
    // No need to render if only one page
    if (totalPages <= 1) return null;

    // Calculate which page numbers to show 
    const getPageNumbers = () => {
        const pages: (number | string)[] = [];
        const showEllipsis = totalPages > 7;

        if (!showEllipsis) {
            for (let i = 1; i <= totalPages; i++) {
                pages.push(i);
            }
        } else {
            // Always show first, last, current, and neighbors
            pages.push(1);

            if (currentPage <= 3) {
                pages.push(2, 3, 4, '...', totalPages);
            } else if (currentPage >= totalPages - 2) {
                pages.push('...', totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
            } else {
                pages.push('...', currentPage - 1, currentPage, currentPage + 1, '...', totalPages);
            }
        }

        return pages;
    }

    const pageNumbers = getPageNumbers();

    return (
        <Box
            sx={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                gap: 1,
                mt: 2,
                mb: 4,
                flexWrap: 'wrap',
            }}
        >
            {/* First Page */}
            <IconButton
                size="small"
                onClick={() => onPageChange(1)}
                disabled={disabled || currentPage === 1}
                title="First Page"
            >
                <FirstPage />
            </IconButton>

            {/* Previous Page */}
            <IconButton
                size="small"
                onClick={() => onPageChange(currentPage - 1)}
                disabled={disabled || currentPage === 1}
                title="Previous Page"
            >
                <ChevronLeft />
            </IconButton>

            {/* Page Numbers */}
            {pageNumbers.map(
                (page, index) => {
                    if (page === '...') {
                        return (
                            <Typography
                                key={`ellipsis-${index}`}
                                variant="body2"
                                sx={{ 
                                    px: 1,
                                    color: 'text.disabled',
                                 }}
                            >
                                ...
                            </Typography>
                        );
                    }

                    const pageNumber = page as number;
                    const isActivePage = pageNumber === currentPage;

                    return (
                        <Button
                            key={pageNumber}
                            variant={isActivePage ? 'contained' : 'outlined'}
                            size="small"
                            disabled={disabled}
                            onClick={() => onPageChange(pageNumber)}
                            sx={{
                                minWidth: 40,
                                height: 40,
                                fontWeight: isActivePage ? 'bold' : 'normal',
                            }}
                        >
                            {pageNumber}
                        </Button>
                    )
                } 
            )}

            {/* Next Page */}
            <IconButton
                size="small"
                onClick={() => onPageChange(currentPage + 1)}
                disabled={disabled || currentPage === totalPages}
                title="Next Page"
            >
                <ChevronRight/>
            </IconButton>

            {/* Last Page */}
            <IconButton
                size="small"
                onClick={() => onPageChange(totalPages)}
                disabled={disabled || currentPage === totalPages}
                title="Last Page"
            >
                <LastPage />
            </IconButton>

            {/* Current Page Info */}
            <Typography
                variant="body2"
                color="text.secondary"
                sx={{ 
                    ml: 2,
                    whiteSpace: 'nowrap',
                }}
            >
                Page {currentPage} of {totalPages}
            </Typography>
        </Box>
    );
}