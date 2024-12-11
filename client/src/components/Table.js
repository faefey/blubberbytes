// Table.js
import React, { useState, useEffect, useMemo } from 'react';
import '../stylesheets/table.css';
import TableContext from './TableContext';
import SelectedFileMenu from './selectedfilemenu';
import checkboxOn from '../icons/checkbox_on.svg';
import checkboxOff from '../icons/checkbox_off.svg';

// TableHead component
const TableHead = ({ columns, handleSorting, selectAll, onSelectAll }) => {
	const [sortField, setSortField] = useState('');
	const [order, setOrder] = useState('asc');

	const handleSortingChange = (accessor) => {
		const sortOrder =
			accessor === sortField && order === 'asc' ? 'desc' : 'asc';
		setSortField(accessor);
		setOrder(sortOrder);
		handleSorting(accessor, sortOrder);
	};

	return (
		<thead>
			<tr>
				<th>
					<img
						src={selectAll ? checkboxOn : checkboxOff}
						alt={selectAll ? 'Select All Checked' : 'Select All Unchecked'}
						className="custom-checkbox"
						onClick={onSelectAll}
					/>
				</th>
				{columns.map(({ label, accessor, sortable }) => {
					const cl = sortable
						? sortField === accessor && order === 'asc'
							? 'up'
							: sortField === accessor && order === 'desc'
								? 'down'
								: 'default'
						: '';
					return (
						<th
							key={accessor}
							onClick={sortable ? () => handleSortingChange(accessor) : null}
							className={cl}
						>
							{label}
						</th>
					);
				})}
			</tr>
		</thead>
	);
};

const TableBody = ({ tableData, columns, onSelectRow, selectedRows }) => {
	return (
		<tbody>
			{tableData.map((data) => {
				const isSelected = selectedRows.some((file) => file.hash === data.hash);
				return (
					<tr
						key={data.hash}
						className={isSelected ? "selected" : ""}
						onClick={() => onSelectRow(data)}
					>
						<td onClick={(e) => e.stopPropagation()}>
							<img
								src={isSelected ? checkboxOn : checkboxOff}
								alt={isSelected ? "Row Selected" : "Row Unselected"}
								className="custom-checkbox"
								onClick={() => onSelectRow(data)}
							/>
						</td>
						{columns.map(({ accessor }) => {
							let tData = data[accessor] !== undefined ? data[accessor] : "——";

							if (accessor === "size") {
								tData = formatSize(tData);
							}

							return <td key={accessor}>{tData}</td>;
						})}
					</tr>
				);
			})}
		</tbody>
	);
};


// Main Table component
const Table = ({ currSection, data, addFile, removeFiles }) => {
	const columns = useMemo(() => {
		switch (currSection) {
			case 'hosting':
				return [
					{ label: 'Hash', accessor: 'hash', sortable: true },
					{ label: 'Name', accessor: 'name', sortable: true },
					{ label: 'Extension', accessor: 'extension', sortable: true },
					{ label: 'Size', accessor: 'size', sortable: true },
					{ label: 'Date', accessor: 'date', sortable: true, sortbyOrder: 'desc' },
					{ label: 'Price', accessor: 'price', sortable: true },
				];
			case 'sharing':
				return [
					{ label: 'Hash', accessor: 'hash', sortable: true },
					{ label: 'Name', accessor: 'name', sortable: true },
					{ label: 'Extension', accessor: 'extension', sortable: true },
					{ label: 'Size', accessor: 'size', sortable: true },
					{ label: 'Date', accessor: 'date', sortable: true, sortbyOrder: 'desc' },
					{ label: 'Password', accessor: 'password', sortable: false },
				];
			case 'saved':
				return [
					{ label: 'Hash', accessor: 'hash', sortable: true },
					{ label: 'Name', accessor: 'name', sortable: true },
					{ label: 'Extension', accessor: 'extension', sortable: true },
					{ label: 'Size', accessor: 'size', sortable: true },

				];
			default:
				return [
					{ label: 'Hash', accessor: 'hash', sortable: true },
					{ label: 'Name', accessor: 'name', sortable: true },
					{ label: 'Extension', accessor: 'extension', sortable: true },
					{ label: 'Size', accessor: 'size', sortable: true },
					{ label: 'Date', accessor: 'date', sortable: true , sortbyOrder: 'desc' },
				];
		}
	}, [currSection]);

	const [filters, setFilters] = useState({
		type: '',
		size: '',
		date: '',
		price: '',
	});

	const [selectedFiles, setSelectedFiles] = useState([]);

	const [tableData, handleSorting] = useSortableTable(data, columns, filters);

	useEffect(() => {
		setSelectedFiles([]);
		setFilters({
			type: '',
			size: '',
			date: '',
			price: '',
		});
	}, [data]);

	const onSelectRow = (file) => {
		setSelectedFiles((prevSelectedFiles) =>
			prevSelectedFiles.some((selectedFile) => selectedFile.hash === file.hash)
				? prevSelectedFiles.filter(
					(selectedFile) => selectedFile.hash !== file.hash
				)
				: [...prevSelectedFiles, file]
		);
	};

	const onSelectAll = () => {
		if (selectedFiles.length === tableData.length) {
			setSelectedFiles([]); // Deselect all
		} else {
			setSelectedFiles([...tableData]); // Select all
		}
	};

	const selectAll =
		selectedFiles.length === tableData.length && tableData.length > 0;

	const contextValue = {
		filters,
		setFilters,
		selectedFiles,
		setSelectedFiles,
	};

	return (
		<TableContext.Provider value={contextValue}>
			<SelectedFileMenu
				currSection={currSection}
				addFile={addFile}
				removeFiles={removeFiles}
			/>
			<table className="table">
				<caption></caption>
				<TableHead
					columns={columns}
					handleSorting={handleSorting}
					selectAll={selectAll}
					onSelectAll={onSelectAll}
				/>
				<TableBody
					columns={columns}
					tableData={tableData}
					onSelectRow={onSelectRow}
					selectedRows={selectedFiles}
				/>
			</table>
		</TableContext.Provider>
	);
};

export default Table;

// Sorting and filtering logic
function getDefaultSorting(defaultTableData, columns) {
	const filterColumn = columns.filter((column) => column.sortbyOrder);

	let { accessor = columns[0]?.accessor || 'hash', sortbyOrder = 'asc' } =
		Object.assign({}, ...filterColumn);

	const sorted = [...defaultTableData].sort((a, b) => {
		const aValue = a[accessor];
		const bValue = b[accessor];

		// Check for undefined or null
		if (aValue == null && bValue == null) return 0;
		if (aValue == null) return 1;
		if (bValue == null) return -1;

		let comparison = 0;

		if (typeof aValue === 'number' && typeof bValue === 'number') {
			comparison = aValue - bValue;
		} else {
			comparison = aValue.toString().localeCompare(bValue.toString(), 'en', {
				numeric: true,
			});
		}

		return sortbyOrder === 'asc' ? comparison : -comparison;
	});
	return sorted;
}

const useSortableTable = (data, columns, filters) => {
	const [tableData, setTableData] = useState([]);

	useEffect(() => {
		let filteredData = applyFilters(data, filters);
		const sortedData = getDefaultSorting(filteredData, columns);
		setTableData(sortedData);
	}, [data, columns, filters]);

	const handleSorting = (sortField, sortOrder) => {
		if (sortField) {
			const sorted = [...tableData].sort((a, b) => {
				const aValue = a[sortField];
				const bValue = b[sortField];

				if (aValue == null && bValue == null) return 0;
				if (aValue == null) return 1;
				if (bValue == null) return -1;

				let comparison = 0;

				if (typeof aValue === 'number' && typeof bValue === 'number') {
					comparison = aValue - bValue;
				} else {
					comparison = aValue.toString().localeCompare(bValue.toString(), 'en', {
						numeric: true,
					});
				}

				return sortOrder === 'asc' ? comparison : -comparison;
			});
			setTableData(sorted);
		}
	};

	return [tableData, handleSorting];
};

function applyFilters(data, filters) {
	return data.filter((item) => {
		let isValid = true;

		// Size filter
		if (filters.size && item.size != null) {
			const bytes = item.size;
			if (filters.size === 'less1mb') {
				// less than 1 MB
				isValid = isValid && bytes < 1e6;
			} else if (filters.size === 'less1gb') {
				// less than 1 GB
				isValid = isValid && bytes < 1e9;
			} else if (filters.size === 'more1gb') {
				// more than 1 GB
				isValid = isValid && bytes > 1e9;
			}
		}

		// Date filter
		if (filters.date && item.date) {
			const itemDate = new Date(item.date);
			const today = new Date();
			if (filters.date === 'today') {
				isValid = isValid && itemDate.toDateString() === today.toDateString();
			} else if (filters.date === '7days') {
				const lastWeek = new Date();
				lastWeek.setDate(today.getDate() - 7);
				isValid = isValid && itemDate >= lastWeek && itemDate <= today;
			} else if (filters.date === '30days') {
				const lastMonth = new Date();
				lastMonth.setDate(today.getDate() - 30);
				isValid = isValid && itemDate >= lastMonth && itemDate <= today;
			} else if (filters.date === '6months') {
				const lastSixMonths = new Date();
				lastSixMonths.setMonth(today.getMonth() - 6);
				isValid = isValid && itemDate >= lastSixMonths && itemDate <= today;
			} else if (filters.date === 'thisyear') {
				isValid = isValid && itemDate.getFullYear() === today.getFullYear();
			} else if (filters.date === 'lastyear') {
				isValid =
					isValid && itemDate.getFullYear() === today.getFullYear() - 1;
			}
		}

		// Price filter (for Hosting)
		// Price filter
		if (filters.price && item.price != null) {
			if (filters.price === 'less5') {
				isValid = isValid && item.price < 5;
			} else if (filters.price === '5to20') {
				isValid = isValid && item.price >= 5 && item.price <= 20;
			} else if (filters.price === 'more20') {
				isValid = isValid && item.price > 20;
			}
		}


		// Type filter
		if (filters.type && item.name) {
			// Extract file extension
			const ext = item.name.includes(".") ? item.name.split(".").pop().toLowerCase() : ""; // Handle no extension

			// Debug print statement

			let fileType = 'other';

			// Define document and media extensions
			const documentExtensions = ['pdf', 'txt', 'doc', 'docx']; 
			const mediaExtensions = ['png', 'jpg', 'jpeg', 'mp4', 'mp3'];

			// Check type based on extension
			if (documentExtensions.includes(ext)) {
				fileType = 'document';
			} else if (mediaExtensions.includes(ext)) {
				fileType = 'media';
			}


			if (filters.type !== fileType) {
				isValid = false;
			}
		}

		

		return isValid;
	});
}




export function formatSize(bytes) {
	if (bytes >= 1e9) {
		return (bytes / 1e9).toFixed(2) + ' GB';
	} else if (bytes >= 1e6) {
		return (bytes / 1e6).toFixed(2) + ' MB';
	} else if (bytes >= 1e3) {
		return (bytes / 1e3).toFixed(2) + ' KB';
	} else {
		return bytes + ' B';
	}
}
