export interface Page<T> {
    items: T[];
    page: number;
    page_size: number;
    total: number;
    has_next_page: boolean;
}