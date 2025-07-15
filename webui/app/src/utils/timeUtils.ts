/**
 * 时间处理工具函数
 * 提供统一的时间格式化和处理方法
 */

/**
 * 格式化日期为本地日期字符串
 * @param dateString - ISO 8601 格式的日期字符串或其他可解析的日期格式
 * @returns 格式化后的本地日期字符串 (YYYY-MM-DD 格式)
 */
export const formatDisplayDate = (dateString: string): string => {
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      console.warn(`Invalid date string: ${dateString}`);
      return dateString; // 如果无法解析，返回原字符串
    }
    return date.toLocaleDateString();
  } catch (error) {
    console.error(`Error formatting date: ${dateString}`, error);
    return dateString;
  }
};

/**
 * 格式化日期为完整的本地日期时间字符串
 * @param dateString - ISO 8601 格式的日期字符串或其他可解析的日期格式
 * @returns 格式化后的本地日期时间字符串
 */
export const formatDisplayDateTime = (dateString: string): string => {
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      console.warn(`Invalid date string: ${dateString}`);
      return dateString;
    }
    return date.toLocaleString();
  } catch (error) {
    console.error(`Error formatting datetime: ${dateString}`, error);
    return dateString;
  }
};

/**
 * 获取用于排序的时间戳
 * @param dateString - ISO 8601 格式的日期字符串或其他可解析的日期格式
 * @returns 时间戳（毫秒）
 */
export const getTimestampForSorting = (dateString: string): number => {
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
      console.warn(`Invalid date string for sorting: ${dateString}`);
      return 0; // 无效日期返回0，会被排在最前面
    }
    return date.getTime();
  } catch (error) {
    console.error(`Error getting timestamp: ${dateString}`, error);
    return 0;
  }
};

/**
 * 验证日期字符串是否有效
 * @param dateString - 要验证的日期字符串
 * @returns 是否为有效日期
 */
export const isValidDate = (dateString: string): boolean => {
  try {
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  } catch {
    return false;
  }
};

/**
 * 为数据对象添加显示日期字段
 * @param data - 包含时间字段的数据对象
 * @param timeField - 时间字段名
 * @returns 添加了displayDate字段的数据对象
 */
export const addDisplayDate = <T extends Record<string, any>>(
  data: T,
  timeField: keyof T = 'created_at'
): T & { displayDate: string } => {
  return {
    ...data,
    displayDate: formatDisplayDate(data[timeField] as string)
  };
};

/**
 * 批量为数据数组添加显示日期字段
 * @param dataArray - 数据对象数组
 * @param timeField - 时间字段名
 * @returns 添加了displayDate字段的数据对象数组
 */
export const addDisplayDateToArray = <T extends Record<string, any>>(
  dataArray: T[],
  timeField: keyof T = 'created_at'
): (T & { displayDate: string })[] => {
  return dataArray.map(item => addDisplayDate(item, timeField));
};