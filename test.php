<?php
// Create a new cURL resource.
$ch = curl_init('http://localhost:8080/');

// Set the HTTP request method to POST.
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, 'POST');

// Set the request headers to include a Content-Type header with a value of application/json.
curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-Type: application/json'));

// Set the request body to the JSON data you want to send.
$data = [
    "from" => "john doe <john@example.com>",
    "to" => ["gong@example.com", "gong@simple.com"],
    "cc" => ["another@example.com", "later@sample.com"],
    "bcc" => ["hidden@example.com", "hiddenalso@example.com"],
    "subject" => "Test email from json2smtp 90",
    "message" => "This is a test email from json2smtp.",
    "attachments" => ["file-sample_1MB.doc"=>base64_encode(file_get_contents('file-sample_1MB.doc'))],
    // "smtphost" => "smtp.example.com",
    // "smtpport" => 587,
    // "smtpuser" => "username",
    // "smtppassword" => "password"
];
$json_data = json_encode($data);
curl_setopt($ch, CURLOPT_POSTFIELDS, $json_data);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

// Execute the request.
$response = curl_exec($ch);

// Check the response status code.
$status_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
if ($status_code !== 200) {
  echo 'Error: ' . curl_error($ch);
  exit;
}

// Handle the response.
echo $response;

// Close the cURL resource.
curl_close($ch);